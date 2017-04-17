package front

type TradeState int

const (
	UNKNOWN TradeState = iota
	NOTPAY
	SUCCESS
	REFUND
	CLOSED
	REVOKED
	USERPAYING
	PAYERROR
)

type OrderState int

const (
	TOrderStateUnknown   OrderState = iota // zh:"未知状态"
	TOrderStateNopay                       // zh:"待付款"
	TOrderStateCanceled                    // zh:"已关闭"
	TOrderStatePrepaid                     // zh:"支付中"
	TOrderStatePaid                        // zh:"已付款"
	TOrderStateRejected                    // zh:"已拒绝"
	TOrderStateEnsuring                    // zh:"等待确认"
	TOrderStateEnsured                     // zh:"已确认"
	TOrderStateRefund                      // zh:"已退款"
	TOrderStateCompleted                   // zh:"已完成"
	TOrderStateEvaled                      // zh:"已评价"
	TOrderStateHistory                     // zh:"已评价"
)

type Order struct {
	tableName struct{} `sql:"cc_order"`
	ID        uint
	UserID    uint `json:"-"` // check owner only
	Remark    string

	PackageID uint
	Unit      string // `zh:"单位"`
	Quantity  uint
	Price     uint
	Name      string

	PayAmount    uint
	WxPaid       uint
	WxRefund     uint
	CashPaid     uint
	CashRefund   uint
	RefundReason string

	// study http://help.vipshop.com/themelist.php?type=detail&id=330
	State       OrderState
	CreatedAt   int64
	CanceledAt  int64
	PrepaidAt   int64 // if by wx
	PaidAt      int64
	RejectedAt  int64 // by operator, trigger auto refound
	EnsuredAt   int64 // by operator if needed
	RefundAt    int64 // by operator
	CompletedAt int64 // resource allocated
	EvalAt      int64 // gen by server
	HistoryAt   int64 // no modify any more

	// EvalItem
	Eval     string
	EvalName string // gen by server
	RateStar uint

	// auto set by system
	NeedEnsure    bool
	AutoCompleted bool
	AutoEvaled    bool

	Rebated bool `json:"-"`
	User1   uint `json:"-"`

	// weixin
	WxPrepayID      string     `json:"-"`
	WxTransactionId string     `json:"-"`
	WxTradeState    TradeState `json:"-"`
	WxRefundID      string     `json:"-"`
	WxTradeNo       string     `json:"-"`
}

type CheckoutPayload struct {
	SkuID    uint
	Quantity uint
	Remark   string
	Total    uint // final amount to pay, used to validate
}

type OrderChangeStatePayload struct {
	ID    uint
	State OrderState
}

type OrderWxPayPayload struct {
	OrderID uint
}

type OrderPrepayResponse struct {
	Order     *Order
	WxPayArgs *WxPayArgs
}

type OrderPayPayload struct {
	Key     string
	OrderID uint
	Amount  uint
}
