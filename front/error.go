package front

import (
	"github.com/empirefox/bongine/cerr"
	"github.com/gin-gonic/gin"
)

type ErrorsParser int

const (
	ParseGinValidator ErrorsParser = iota // github.com/go-playground/validator. TODO rename
	ParseGovalidator                      // github.com/asaskevich/govalidator
	ParseEfValidator
	ParseCode
)

type Errors struct {
	Parser ErrorsParser
	Code   cerr.CodedError
	Err    interface{}
}

func NewErrors(parser ErrorsParser, err interface{}) *Errors {
	return &Errors{
		Parser: parser,
		Err:    err,
	}
}

func NewGinv(err interface{}) *Errors {
	return &Errors{
		Parser: ParseGinValidator,
		Err:    err,
	}
}

func NewGov(err interface{}) *Errors {
	return &Errors{
		Parser: ParseGovalidator,
		Err:    err,
	}
}

func NewEfv(err interface{}) *Errors {
	return &Errors{
		Parser: ParseEfValidator,
		Err:    err,
	}
}

func NewCodev(code cerr.CodedError) *Errors {
	return &Errors{
		Parser: ParseCode,
		Code:   code,
	}
}

func NewCodeErrv(code cerr.CodedError, err interface{}) *Errors {
	return &Errors{
		Parser: ParseCode,
		Code:   code,
		Err:    err,
	}
}

func (err *Errors) Abort(c *gin.Context, status int) {
	c.JSON(status, err)
	c.Abort()
}

func (err *Errors) AbortIf(c *gin.Context, status int) bool {
	if err == nil {
		return false
	}
	err.Abort(c, status)
	return true
}
