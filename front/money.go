package front

type UserCashType int

const (
	TUserCashUnknown     UserCashType = iota // `zh:"未知"`
	TUserCashPrepay                          // `zh:"预付款"`
	TUserCashPrepayBack                      // `zh:"预付款退回"`
	TUserCashTrade                           // `zh:"交易"`
	TUserCashRefund                          // `zh:"退款"`
	TUserCashPreWithdraw                     // `zh:"预提现"`
	TUserCashWithdraw                        // `zh:"提现"`
	TUserCashReward                          // `zh:"奖励"`
	TUserCashRebate                          // `zh:"返利"`
	TUserCashRecharge                        // `zh:"充值"`
)

type UserCash struct {
	tableName struct{} `sql:"cc_user_cash,alias:uc"`
	ID        uint
	UserID    uint `json:"-"`
	OrderID   uint
	CreatedAt int64
	Type      UserCashType
	Amount    int
	Remark    string
	Balance   int
}

type UserCashFrozen struct {
	tableName struct{} `sql:"cc_user_cash_frozen,alias:ucf"`
	ID        uint
	UserID    uint `json:"-"`
	OrderID   uint
	CreatedAt int64
	Type      UserCashType
	Amount    uint
	Remark    string
	ThawedAt  int64
}

type UserCashRebateItem struct {
	tableName struct{} `sql:"cc_user_cash_rebate_item,alias:ucri"`
	ID        uint
	RebateID  uint
	CreatedAt int64
	Amount    uint
}

type UserCashRebate struct {
	tableName struct{} `sql:"cc_user_cash_rebate,alias:ucr"`
	ID        uint
	UserID    uint `json:"-"`
	OrderID   uint
	CreatedAt int64
	Type      UserCashType
	Amount    uint
	Remark    string
	Stages    uint
	DoneAt    int64

	Items []UserCashRebateItem `pg:",fk:Rebate"`
}

// cash but without withdraw
type PointsItem struct {
	tableName struct{} `sql:"cc_points_item,alias:pi"`
	ID        uint
	UserID    uint `json:"-"`
	TaskID    uint
	CreatedAt int64
	Amount    int
	Balance   int
}

type WxPayArgs struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type WithdrawPayload struct {
	Amount uint
	Ip     string `json:"-"`
}
