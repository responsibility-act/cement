package ferr

import (
	"fmt"

	"github.com/empirefox/cement/perr"
)

var EfFieldErrorType = 101

type EfFieldError struct {
	TagType int
	Tag     int
	Field   int
	Param   int
}

func (e EfFieldError) Coded() *CodedErr {
	return &perr.CodedErr{T: EfFieldErrorType, C: int(e)}
}

func (e EfFieldError) Error() string {
	return fmt.Sprintf("EfFieldError type: %d, code: %d", e.TagType, e)
}

func init() {
	perr.Check(EfFieldError{TagType: EfFieldErrorType})
}
