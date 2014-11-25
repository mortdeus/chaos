package chaos

import (
	"fmt"
	"reflect"
	"strings"
)

func typeCheck(sval1, sval2, dval reflect.Value) (reflect.Kind, error) {

	if sval1.Type().AssignableTo(sval2.Type()) {
		if sval1.Type().AssignableTo(dval.Elem().Type()) {
			return sval1.Kind(), nil
		}
	}
	return reflect.Invalid, TypeErr{fmt.Errorf("Incompatible types: \"%v = %v + %v\"", dval, sval1, sval2)}
}

type TypeErr struct {
	error
}

const (
	Unknown = iota - 1
	InvalidOp
	TypeConv
)

func (t TypeErr) Code() int {
	switch strings.Split(t.error.Error(), ":")[0] {
	case "Incompatible types":
		return TypeConv
	case "Invalid operation":
		return InvalidOp
	default:
		return Unknown
	}
}
