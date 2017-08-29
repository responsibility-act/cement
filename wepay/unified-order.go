package wepay

import (
	"strconv"
	"time"

	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/dchest/uniuri"
	"github.com/empirefox/cement/perr"
	"go.uber.org/zap"
)

func (wc *Wepay) UnifiedOrder(order *UnifiedOrder) (*UnifiedOrderResponse, error) {
	senario, ok := Senarios[order.Senario]
	if !ok {
		wc.log.Error("Senario error", zap.String("senario", order.Senario.String()))
		return nil, perr.FrontendErr
	}
	req := &pay.UnifiedOrderRequest{
		DeviceInfo:     senario.DeviceInfo,
		Body:           wc.config.PayBody,
		OutTradeNo:     order.OutTradeNo,
		TotalFee:       int64(order.TotalFee),
		FeeType:        order.FeeType,
		SpbillCreateIP: order.Ip,
		NotifyURL:      wc.config.PayNotifyURL,
		TradeType:      senario.TradeType,
		OpenId:         order.OpenId,
	}

	res, err := pay.UnifiedOrder2(wc.Client, req)
	if err != nil {
		if err2, ok := err.(*core.BizError); ok {
			switch err2.ErrCode {
			case "NOTENOUGH", "ORDERPAID", "ORDERCLOSED", "SYSTEMERROR", "OUT_TRADE_NO_USED":
				// not source code error, no remote log
				wc.log.Info("result FAIL", zap.String("err_code", err2.ErrCode))
				return nil, perr.WxResultCodeFail
			default:
				// log remote
				wc.log.Error("result FAIL", zap.String("err_code", err2.ErrCode))
				return nil, perr.WxResultCodeFail
			}
		}
		wc.log.Error("return/result FAIL", zap.Error(err))
		return nil, perr.WxReturnCodeFail
	}

	switch order.Senario {
	case TSenarioWechat:
		return &UnifiedOrderResponse{Wechat: wc.JsapiSign(&res.PrepayId)}, nil
	case TSenarioH5:
		return &UnifiedOrderResponse{H5: &res.MWebURL}, nil
	case TSenarioNative2:
		return &UnifiedOrderResponse{Native2: &res.CodeURL}, nil
	}
	wc.log.Error("Senario interal error", zap.String("senario", order.Senario.String()))
	return nil, perr.InternalErr
}

func (wc *Wepay) JsapiSign(prepayId *string) *JsapiArgs {
	args := &JsapiArgs{
		AppId:     wc.config.AppId,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  uniuri.NewLen(32),
		Package:   "prepay_id=" + *prepayId, // 2hour
		SignType:  "MD5",
	}
	args.PaySign = core.JsapiSign(args.AppId, args.TimeStamp, args.NonceStr, args.Package, args.SignType, wc.config.MchKey)
	return args
}
