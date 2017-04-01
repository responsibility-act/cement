package wx

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
)

func (wc *WxClient) SendCoupon(toUsr, stockId string, id uint) {
	req := map[string]string{
		"coupon_stock_id":  stockId,
		"openid_count":     "1",
		"partner_trade_no": fmt.Sprintf("%s%s%d", wc.wx.MchId, time.Now().Format("20060102"), id),
		"openid":           toUsr,
		"appid":            wc.wx.AppId,
		"mch_id":           wc.wx.MchId,
		"nonce_str":        uniuri.NewLen(32),
	}
	_ = req
}
