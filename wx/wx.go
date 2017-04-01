package config

import (
	"crypto/md5"
	"encoding/xml"
	"io"
	"strconv"
	"time"

	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat.v2/mch/core"
	"github.com/chanxuehong/wechat.v2/mch/mmpaymkttransfers/promotion"
	"github.com/chanxuehong/wechat.v2/mch/pay"
	"github.com/dchest/uniuri"
	"github.com/empirefox/bongine/config"
	"github.com/empirefox/bongine/front"
	"github.com/empirefox/bongine/models"
)

type WxClient struct {
	*core.Client
	hclient   *core.Client
	notifyUrl string
	config    config.Weixin
}

func NewWxClient(config *config.Config) (*WxClient, error) {
	weixin := &config.Weixin
	httpClient, err := core.NewTLSHttpClient(weixin.CertFile, weixin.KeyFile)
	if err != nil {
		return nil, err
	}
	return &WxClient{
		config:    *weixin,
		notifyUrl: config.Security.SecendOrigin + weixin.PayNotifyURL,
		Client:    core.NewClient(weixin.AppId, weixin.MchId, weixin.MchKey, httpClient),
		hclient:   core.NewClient(weixin.AppId, weixin.MchId, weixin.MchKey, nil),
	}, nil
}

func (wc *WxClient) NewWxPayArgs(prepayId *string) *front.WxPayArgs {
	args := &front.WxPayArgs{
		AppId:     wc.config.AppId,
		TimeStamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  uniuri.NewLen(32),
		Package:   "prepay_id=" + *prepayId, // 2hour
		SignType:  "MD5",
	}
	args.PaySign = core.JsapiSign(args.AppId, args.TimeStamp, args.NonceStr, args.Package, args.SignType, wc.config.MchKey)
	return args
}

// only can be called by PrepayOrder
func (wc *WxClient) UnifiedOrder(tokUsr *models.User, order *front.Order, ip *string) (*string, *front.WxPayArgs, error) {
	req := &pay.UnifiedOrderRequest{
		DeviceInfo:     "WEB",
		Body:           wc.config.PayBody,
		OutTradeNo:     order.WxOutTradeNo(),
		TotalFee:       int64(order.PayAmount),
		SpbillCreateIP: *ip,
		NotifyURL:      wc.notifyUrl,
		TradeType:      "JSAPI",
		OpenId:         tokUsr.OpenId,
	}

	res, err := pay.UnifiedOrder2(wc.Client, req)
	if err != nil {
		return nil, nil, err
	}

	return &res.PrepayId, wc.NewWxPayArgs(&res.PrepayId), nil
}

func (wc *WxClient) OnWxPayNotify(r io.Reader) (*WxResponse, map[string]string) {
	m, err := util.DecodeXMLToMap(r)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse request body"), nil
	}
	if m["return_code"] != "SUCCESS" {
		return NewWxResponse(m["return_code"], m["return_msg"]), nil
	}

	sign := core.Sign(m, wc.config.MchKey, md5.New)
	if sign != m["sign"] {
		return NewWxResponse("FAIL", "failed to validate md5"), nil
	}

	if m["result_code"] == "SUCCESS" {
		m["trade_state"] = "SUCCESS"
	} else {
		m["trade_state"] = "NOPAY"
	}

	return nil, m
}

func (wc *WxClient) OrderQuery(order *front.Order) (map[string]string, error) {
	req := map[string]string{
		"appid":        wc.config.AppId,
		"mch_id":       wc.config.MchId,
		"out_trade_no": order.WxOutTradeNo(),
		"nonce_str":    uniuri.NewLen(32),
	}
	if order.WxTransactionId != "" {
		req["transaction_id"] = order.WxTransactionId
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	return pay.OrderQuery(wc.Client, req)
}

func (wc *WxClient) OrderClose(order *front.Order) (map[string]string, error) {
	req := map[string]string{
		"appid":        wc.config.AppId,
		"mch_id":       wc.config.MchId,
		"out_trade_no": order.WxOutTradeNo(),
		"nonce_str":    uniuri.NewLen(32),
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	return pay.CloseOrder(wc.Client, req)
}

func (wc *WxClient) OrderRefund(order *front.Order) (map[string]string, error) {
	req := map[string]string{
		"appid":         wc.config.AppId,
		"mch_id":        wc.config.MchId,
		"nonce_str":     uniuri.NewLen(32),
		"out_trade_no":  order.WxOutTradeNo(),
		"out_refund_no": order.WxOutTradeNo(),
		"total_fee":     strconv.Itoa(int(order.WxPaid)),
		"refund_fee":    strconv.Itoa(int(order.WxRefund)),
		"op_user_id":    wc.config.MchId,
	}
	req["sign"] = core.Sign(req, wc.config.MchKey, md5.New)
	return pay.Refund(wc.Client, req)
}

type TransfersArgs struct {
	TradeNo string
	OpenID  string
	Amount  uint
	Desc    string
	Ip      string
}

func (wc *WxClient) Transfers(args *TransfersArgs) (map[string]string, error) {
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
	return promotion.Transfers(wc.Client, req)
}

type WxResponse struct {
	XMLName    xml.Name  `xml:"xml"`
	ReturnCode CDATAText `xml:"return_code"`
	ReturnMsg  CDATAText `xml:"return_msg,omitempty"`
}

type CDATAText struct {
	Text string `xml:",cdata"`
}

func NewWxResponse(code, msg string) *WxResponse {
	return &WxResponse{
		ReturnCode: CDATAText{code},
		ReturnMsg:  CDATAText{msg},
	}
}
