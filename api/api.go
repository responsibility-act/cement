package api

import "github.com/mcuadros/go-defaults"

type Apis struct {
	GetProfile         string `default:"/profile"`
	GetUserResources   string `default:"/resources"`
	GetUserResource    string `default:"/resource/id/:id"`
	PatchResourceMount string `default:"/resource/mount"`
	GetUserSites       string `default:"/sites"`
	GetUserSite        string `default:"/site/:id"`
	PostUserSite       string `default:"/site"`
	DeleteUserSite     string `default:"/site/:id"`
	GetCaptcha         string `default:"/captcha"`
	GetRefreshToken    string `default:"/refresh_token/:refreshToken"`
	GetFakeToken       string `default:"/faketoken"`
	GetFans            string `default:"/fans"`
	PostSetUserInfo    string `default:"/user"`
	PostPrebindPhone   string `default:"/phone/prebind"`
	PostBindPhone      string `default:"/phone/bind"`
	PostPresetPaykey   string `default:"/paykey/preset"`
	PostSetPaykey      string `default:"/paykey/set"`
	PostWithdraw       string `default:"/withdraw"`
	GetUserCash        string `default:"/cash"`
	GetUserCashFrozen  string `default:"/cash/frozen"`
	GetUserCashRebate  string `default:"/cash/rebate"`
	GetUserPoints      string `default:"/points"`
	GetOrder           string `default:"/order/:id"`
	GetOrders          string `default:"/orders"`
	PostCheckout       string `default:"/checkout"`

	PostPayCash       string `default:"/pay/cash"`
	PostPayPoints     string `default:"/pay/points"`
	PostWepayInWechat string `default:"/wepay/wechat"`
	PostWepayInH5     string `default:"/wepay/h5"`
	PostWepayInWithQr string `default:"/wepay/qr"`
	PostWepayAfterPay string `default:"/wepay/paid"`

	PostOrderState string `default:"/order/state"`
	PostOrderEval  string `default:"/order/eval/:id"`
	GetPackages    string `default:"/packages"`

	GetQiniuCommon string `default:"/qiniu/commons"`

	// GetQiniuHeadToken only upload with key: h/[userid]. Token should be got early.
	GetQiniuHeadToken string `default:"/qiniu/headtoken/:life"`

	// GetQiniuUptoken only upload with key. It will check site's owner.
	// Called when upload action indeed happens.
	// name = base64('s/:siteid/2017/01/abc.png')
	GetQiniuUptoken string `default:"/qiniu/uptoken/:key/:life"`

	// PostQiniuList only list with prefix. It will check site's owner.
	// prefix = base64('s/:siteid/2017/01/')
	PostQiniuList string `default:"/qiniu/:prefix"`

	// GetQiniuList only delete with key. It will check site's owner.
	// name = base64('s/:siteid/2017/01/abc.png')
	DeleteQiniu          string `default:"/qiniu/:key"`
	PostBatchDeleteQiniu string `default:"/qiniu/batch"`

	PostAuthWx string `default:"/auth/wx"` // weixin only
}

func NewApis() Apis {
	apis := new(Apis)
	defaults.SetDefaults(apis)
	return *apis
}
