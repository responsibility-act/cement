package perr

import (
	"fmt"
	"reflect"
)

type Cerror interface {
	Coded() *CodedErr
	Error() string
}

type CodedErr struct {
	T int
	C int
}

var codeTypes = make(map[int]struct{})

func Check(v Cerror) {
	err := v.Coded()
	_, ok := codeTypes[err.T]
	if ok {
		panic(fmt.Errorf("Error type(%s) exist: %d", reflect.TypeOf(v).String(), err.T))
	}

	codeTypes[err.T] = struct{}{}
}
