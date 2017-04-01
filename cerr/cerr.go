package cerr

import "fmt"

type CodedError int

func (ce CodedError) Error() string {
	return fmt.Sprintf("Error code: %d", ce)
}

// all are non-StatusForbidden
const (
	Error CodedError = iota
	Unauthorized
	SonyFlakeTimeout
	InvalidUrlParam
	InvalidPostBody
	InvalidPhoneFormat
	RebindSamePhone
	PhoneOccupied
	PhoneBindRequired
	UserNotFound
	DbFailed
	CaptchaRejected
	UpdateWxOrderStateFailed
	SystemModeNotAllowed
	InvalidRefreshToken
	NoRefreshToken
	NoNeedRefreshToken
	NoAccessToken
	InvalidAccessToken
	GenCaptchaFailed
	RemoteHTTPFailed
	InvalidTokenSubject
	InvalidTokenExpires
	InvalidSignAlg
	InvalidClaimId
	InvalidUserID
	Forbidden
	RetrySmsFailed
	SendSmsError
	SendSmsFailed
	SmsVerifyFailed

	InvalidProductId
	InvalidSkuStock
	InvalidSkuId
	InvalidAttrId
	InvalidAttrLen
	InvalidGroupbuyId
	InvalidCheckoutTotal
	InvalidCheckoutFreight
	InvalidPaykey
	PaykeyNeedBeSet
	NotEnoughMoney
	NotEnoughPoints
	OnlyAbcOrPoints
	NoAbcOrPoints
	InvalidPayAmount
	InvalidPayType
	NotNopayState
	NoWayToPaidState
	NoWayToTargetState
	NoPermToState
	OrderClosed
	OrderCloseNeeded
	OrderCompleteTimeout
	OrderEvalTimeout
	OrderItemNotFound
	NotPrepayOrder
	WxPayNotCompleted
	WxRefundNotCompleted
	WxOrderNotExist
	WxOrderAlreadyClosed
	WxOrderCloseFailed
	WxOrderAlreadyPaid
	WxOrderCloseIn5Min
	WxSystemFailed
	InvalidCashPrepaid
	InvalidPrepayPayload
	ApiImplementFailed
	WxUserAbnormal
	ParseWxTotalFeeFailed
	WithdrawFailed

	VipRebateSubIDsLen
	VipRebateSubIDsHas0
	VipRebateSubIDsNoRow
	NotVip
	VipBalanceEmpty
	VipRebateSubTotalSmall
	InvalidRebateType
	AmountLimit
)
