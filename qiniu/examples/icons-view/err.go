package main

import (
	"fmt"

	"github.com/empirefox/cement/perr"
)

var ServerErrorType = 20

type ServerError int

func (e ServerError) Coded() *perr.CodedErr {
	return &perr.CodedErr{T: ServerErrorType, C: int(e)}
}

func (e ServerError) Error() string {
	return fmt.Sprintf("Error type: %d, code: %d", ServerErrorType, e)
}

const (
	// BadRequest
	ErrBadParamKey ServerError = iota
	ErrBadParamPrefix
)
