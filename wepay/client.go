package wepay

import (
	"crypto/md5"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers/promotion"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/dchest/uniuri"
	"github.com/empirefox/cement/clog"
	"github.com/empirefox/cement/httphack"
	"github.com/empirefox/cement/perr"
)

type Config struct {
	DevHackHost    string `                           env:"WEPAY_DEV_HACK_HOST"`
	Dev            bool   `validate:"dep=DevHackHost" env:"WEPAY_DEV"`
	AppId          string `validate:"required"`
	MchKey         string `validate:"required"`
	MchId          string `validate:"required"`
	Cert           []byte `validate:"gt=0" json:"-" toml:"-" yaml:"-" xps:"wxapi.crt"`
	Key            []byte `validate:"gt=0" json:"-" toml:"-" yaml:"-" xps:"wxapi.key"`
	PayBody        string `validate:"required"`     // https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_2
	PayNotifyURL   string `validate:"required,url"` // full url
	TransCheckName string
}

type Wepay struct {
	*core.Client
	config Config
	log    *zap.Logger
}

func NewWepay(config Config, logger clog.Logger) (*Wepay, error) {
	lm := logger.Module("wepay")
	var client *http.Client
	var err error
	if config.Dev {
		client = httphack.NewHackClient(config.DevHackHost)
	} else {
		client, err = core.NewTLSHttpClient(string(config.Cert), string(config.Key))
	}

	if err != nil {
		lm.Error("Failed to create http.Client", zap.Error(err))
		return nil, err
	}

	return &Wepay{
		config: config,
		log:    lm,
		Client: core.NewClient(config.AppId, config.MchId, config.MchKey, client),
	}, nil
}

// OrderQuery query order, but need check order like OnNotifyFunc:
// out_trade_no=OutTradeNo total_fee=TotalFee fee_type=FeeType.
// Save trade_state, time_end, transaction_id to database if SUCCESS.
// Common fields with OnNotifyFunc: device_info openid is_subscribe trade_type
// trade_state bank_type total_fee settlement_total_fee fee_type cash_fee
// cash_fee_type coupon_fee coupon_count coupon_type_$n coupon_id_$n
// coupon_fee_$n transaction_id out_trade_no attach time_end
func (wc *Wepay) OrderQuery(order *QueryOrder) (map[string]string, error) {
	req := map[string]string{
		"appid":        wc.config.AppId,
		"mch_id":       wc.config.MchId,
		"out_trade_no": order.OutTradeNo,
		"nonce_str":    uniuri.NewLen(32),
	}
	if order.TransactionId != "" {
		req["transaction_id"] = order.TransactionId
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	res, err := pay.OrderQuery(wc.Client, req)
	if err != nil {
		wc.log.Error("return_code FAIL", zap.Error(err))
		return nil, perr.WxReturnCodeFail
	}
	return res, nil
}

func (wc *Wepay) OrderClose(outTradeNo string) error {
	req := map[string]string{
		"appid":        wc.config.AppId,
		"mch_id":       wc.config.MchId,
		"out_trade_no": outTradeNo,
		"nonce_str":    uniuri.NewLen(32),
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	res, err := pay.CloseOrder(wc.Client, req)
	if err != nil {
		wc.log.Error("return_code FAIL", zap.Error(err))
		return perr.WxReturnCodeFail
	}

	if res["result_code"] != "SUCCESS" {
		switch res["err_code"] {
		case "ORDERPAID":
			err = perr.WxOrderAlreadyPaid
		case "SYSTEMERROR":
			err = perr.WxNeedRetry
		case "ORDERCLOSED": // used as SUCCESS, pass
		default:
			wc.log.Error("result FAIL", zap.String("err_code", res["err_code"]))
			err = perr.WxApiBadImplement
		}
	}
	return err
}

func (wc *Wepay) OrderRefund(order *RefundOrder) error {
	req := map[string]string{
		"appid":           wc.config.AppId,
		"mch_id":          wc.config.MchId,
		"nonce_str":       uniuri.NewLen(32),
		"out_trade_no":    order.OutTradeNo,
		"out_refund_no":   order.OutTradeNo,
		"total_fee":       strconv.Itoa(int(order.TotalFee)),
		"refund_fee":      strconv.Itoa(int(order.RefundFee)),
		"refund_fee_type": order.FeeType,
		"op_user_id":      wc.config.MchId,
	}
	if order.TransactionId != "" {
		req["transaction_id"] = order.TransactionId
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	res, err := pay.Refund(wc.Client, req)
	if err != nil {
		wc.log.Error("return_code FAIL", zap.Error(err))
		return perr.WxReturnCodeFail
	}

	if res["result_code"] != "SUCCESS" {
		switch res["err_code"] {
		case "SYSTEMERROR":
			err = perr.WxNeedRetry
		case "TRADE_OVERDUE":
			wc.log.Error("result FAIL", zap.String("err_code", res["err_code"]))
			err = perr.FrontendErr
		case "ERROR":
			err = perr.WxUnknownErr
		case "USER_ACCOUNT_ABNORMAL":
			err = perr.WxUserAbnormal
		case "INVALID_REQ_TOO_MUCH", "FREQUENCY_LIMITED":
			wc.log.Error("result FAIL", zap.String("err_code", res["err_code"]))
			err = perr.WxNeedCold
		case "NOTENOUGH":
			err = perr.WxNotEnough
		default:
			wc.log.Error("result FAIL", zap.String("err_code", res["err_code"]))
			err = perr.WxApiBadImplement
		}
	}

	return err
}

type TransfersRequest struct {
	TradeNo string
	OpenID  string
	Amount  uint
	Desc    string
	Ip      string
}

func (wc *Wepay) Transfers(args *TransfersRequest) error {
	req := map[string]string{
		"mch_appid":        wc.config.AppId,
		"mchid":            wc.config.MchId,
		"nonce_str":        uniuri.NewLen(32),
		"partner_trade_no": args.TradeNo,
		"openid":           args.OpenID,
		"check_name":       "NO_CHECK",
		"amount":           strconv.FormatUint(uint64(args.Amount), 10),
		"desc":             args.Desc,
		"spbill_create_ip": args.Ip,
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	res, err := promotion.Transfers(wc.Client, req)
	// IMPORTANT!!!
	// There is no sign in response, so ignore it.
	if err == core.ErrNotFoundSign {
		err = nil
	}
	if err != nil {
		wc.log.Error("return_code FAIL", zap.Error(err))
		return perr.WxReturnCodeFail
	}

	if res["result_code"] != "SUCCESS" {
		switch res["err_code"] {
		case "NOTENOUGH":
			err = perr.WxNotEnough
		case "SYSTEMERROR":
			err = perr.WxNeedRetry
		case "V2_ACCOUNT_SIMPLE_BAN":
			err = perr.WxNeedRealPersonAccount
		default:
			// not supported: NAME_MISMATCH
			// api: NOAUTH AMOUNT_LIMIT PARAM_ERROR OPENID_ERROR SIGN_ERROR XML_ERROR FATAL_ERROR CA_ERROR
			wc.log.Error("result FAIL", zap.String("err_code", res["err_code"]))
			err = perr.WxApiBadImplement
		}
	}

	return err
}
