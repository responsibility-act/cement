package wepay

import (
	"github.com/labstack/echo"
)

type WepayPayload struct {
	Senario SenarioType
	OrderID uint
}

// PostPay handle WxService request of ng-ef-sand/pay.
// Wechat: response jsapi args.
// H5: response redirect url to mobile browser.
// Native2: response code url to browser.
func (wc *Wepay) PostPay(c echo.Context) error {
	return nil
}

// PostAfterPay query order from database, then validate it,
// then query wepay server if needed, then save state if needed,
// then response bussiness order to user.
func (wc *Wepay) PostAfterPay(c echo.Context) error {
	return nil
}
