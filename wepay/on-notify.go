package wepay

import (
	"crypto/md5"
	"encoding/xml"
	"net/http"

	"github.com/chanxuehong/util"
	"github.com/chanxuehong/wechat.v2/mch/core"
)

// OnNotifyFunc check out_trade_no=OutTradeNo total_fee=TotalFee fee_type=FeeType.
// Save trade_state, time_end, transaction_id to database if SUCCESS.
// Common fields with OrderQuery: device_info openid is_subscribe trade_type
// trade_state bank_type total_fee settlement_total_fee fee_type cash_fee
// cash_fee_type coupon_fee coupon_count coupon_type_$n coupon_id_$n
// coupon_fee_$n transaction_id out_trade_no attach time_end
type OnNotifyFunc func(res map[string]string) (err error)

func (wc *Wepay) OnNotify(req *http.Request, onNotify OnNotifyFunc, response func(code int, res interface{})) {
	response(http.StatusOK, wc.onNotify(req, onNotify))
}

func (wc *Wepay) onNotify(req *http.Request, onNotify OnNotifyFunc) *WxResponse {
	defer req.Body.Close()
	m, err := util.DecodeXMLToMap(req.Body)
	if err != nil {
		return NewWxResponse("FAIL", "failed to parse request body")
	}
	if m["return_code"] != "SUCCESS" {
		return NewWxResponse(m["return_code"], m["return_msg"])
	}

	sign := core.Sign(m, wc.config.MchKey, md5.New)
	if sign != m["sign"] {
		return NewWxResponse("FAIL", "failed to validate md5")
	}

	if m["result_code"] == "SUCCESS" {
		m["trade_state"] = "SUCCESS"
	} else {
		m["trade_state"] = "NOPAY"
	}

	if err := onNotify(m); err != nil {
		return NewWxResponse("FAIL", err.Error())
	}
	return NewWxResponse("SUCCESS", "")
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
