package perr

import (
	"fmt"
)

var CementErrorType = 100

type CementError int

func (e CementError) Coded() *CodedErr {
	return &CodedErr{T: CementErrorType, C: int(e)}
}

func (e CementError) ToParser() *ParserError {
	return &ParserError{
		Parser: CodedParser,
		Code:   int(e),
	}
}

func (e CementError) Error() string {
	return fmt.Sprintf("Error type: %d, code: %d", CementErrorType, e)
}

const (
	UnknownError CementError = iota
	InternalErr
	FrontendErr
	TryAgain

	// wepay
	WxOrderAlreadyPaid
	WxNeedRetry
	WxApiBadImplement
	WxReturnCodeFail
	WxResultCodeFail
	WxUnknownErr
	WxUserAbnormal
	WxNeedCold
	WxNotEnough
	WxNeedRealPersonAccount

	// sms
	SmsSendFail
	SmsNeedCold

	// captchar
	BadFontPaths

	// qiniu
	QiniuListFail
	QiniuDelFail
)

func init() {
	Check(UnknownError)
}
