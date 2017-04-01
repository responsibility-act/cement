package front

import (
	"time"

	"gopkg.in/pg.v5/orm"
)

type Profile struct {
	tableName            struct{} `sql:"sys_profile"`
	ID                   uint     // always 1
	WxMpName             string
	Phone                string
	Email                string
	DefaultHeadImage     string
	UserCashRebateStages uint `default:"3"`
	CreatedAt            int64
}

func (b *Profile) BeforeInsert(db orm.DB) error {
	b.CreatedAt = time.Now().Unix()
	return nil
}

type ProfileResponse struct {
	*Profile
	WxAppId     string
	WxScope     string
	WxLoginPath string

	// Config.Order
	EvalTimeoutDay        uint
	CompleteTimeoutDay    uint
	HistoryTimeoutDay     uint
	CheckoutExpiresMinute time.Duration
	WxPayExpiresMinute    time.Duration

	HeadPrefix string
}
