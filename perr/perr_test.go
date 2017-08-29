package perr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var MyErrorType = 20

type MyError int

func (e MyError) Coded() *CodedErr {
	return &CodedErr{T: MyErrorType, C: int(e)}
}

func (e MyError) Error() string {
	return fmt.Sprintf("Error type: %d, code: %d", MyErrorType, e)
}

const (
	MyErrorUnknown         MyError = iota // `en:"Unknown error"` for nodejs npt go-consts-ts
	MyErrorUnexpected                     // `en:"Unexpected error"`
	MyErrorInvalidUrlParam                // `en:"Invalid url param"`
	MyErrorInvalidPostBody                // `en:"Invalid post body"`
)

func Test_Check(t *testing.T) {
	require := require.New(t)
	require.NotPanics(func() { Check(MyErrorUnknown) })
}

func Test_NewCode(t *testing.T) {
	require := require.New(t)
	pe := NewCode(MyErrorInvalidPostBody)
	require.Equal(CodedParser, pe.Parser)
	require.Equal(MyErrorType, pe.Code.T)
	require.Equal(int(MyErrorInvalidPostBody), pe.Code.C)
}
