package wepay

type UnifiedOrder struct {
	Senario    SenarioType
	PayBody    string // overwrite config.PayBody https://pay.weixin.qq.com/wiki/doc/api/native.php?chapter=4_2
	OpenId     string
	OutTradeNo string
	TotalFee   int64
	FeeType    string // CNY
	Ip         string
}

type Senario struct {
	DeviceInfo string
	TradeType  string
}

var Senarios = map[SenarioType]Senario{
	TSenarioWechat:  Senario{"WEB", "JSAPI"},
	TSenarioH5:      Senario{"WEB", "MWEB"},
	TSenarioNative2: Senario{"WEB", "NATIVE"},
}

type JsapiArgs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type UnifiedOrderResponse struct {
	Wechat  *JsapiArgs
	H5      *string // MWebURL
	Native2 *string // CodeURL, 2 hour
}

type QueryOrder struct {
	OutTradeNo    string
	TotalFee      int64
	TransactionId string
	FeeType       string // CNY
}

type RefundOrder struct {
	OutTradeNo    string
	TotalFee      int64
	RefundFee     int64
	TransactionId string
	FeeType       string // CNY
}
