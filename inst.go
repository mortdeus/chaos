package chaos

import (
	"fmt"
	"reflect"
)

type Instruction func(src1, src2, dest interface{}) error

func Process(inst Instruction, src1, src2, dest interface{}) error {
	sval1, sval2, dval := reflect.ValueOf(src1), reflect.ValueOf(src2), reflect.ValueOf(dest)
	fmt.Println("inst.go:12  sval1,sval2,dval:", sval1, sval2, reflect.Indirect(dval))

	kind, err := typeCheck(sval1, sval2, dval)
	if err != nil {
		return err
	}

	// Here we do a fancy trick with reflection that allows us to pass
	// in an empty *interface{} as our dest argument and the VM will automatically
	// mutate the interface's internal state to be type compatible with src1's & src2's.

	if dval.Elem().Kind() == reflect.Interface {
		ptr := reflect.New(sval1.Type())
		dval = dval.Elem()
		dval.Set(ptr)
		dval = dval.Elem()
		defer func() {
			dval := reflect.ValueOf(dest)
			i := reflect.Indirect(dval.Elem().Elem().Convert(reflect.PtrTo(sval1.Type()))).Interface()
			reflect.Indirect(dval).Set(reflect.ValueOf(i))
		}()
	}

	dval = dval.Elem()
	var s1, s2, dfn interface{}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		s1, s2, dfn = sval1.Int(), sval2.Int(), dval.SetInt
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		s1, s2, dfn = sval1.Uint(), sval2.Uint(), dval.SetUint
	case reflect.String:
		s1, s2, dfn = sval1.String(), sval2.String(), dval.SetString
	case reflect.Float32, reflect.Float64:
		s1, s2, dfn = sval1.Float(), sval2.Float(), dval.SetFloat
	case reflect.Complex64, reflect.Complex128:
		s1, s2, dfn = sval1.Complex(), sval2.Complex(), dval.SetComplex

	case reflect.Bool, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
		s1, s2, dfn = sval1, sval2, nil
	}
	if err := inst(s1, s2, dfn); err != nil {
		return err
	}
	return nil
}

func Add(src1, src2, dest interface{}) error {
	kind := reflect.ValueOf(src1).Kind()
	switch f := dest.(type) {
	case func(int64):
		f(src1.(int64) + src2.(int64))
	case func(uint64):
		f(src1.(uint64) + src2.(uint64))
	case func(string):
		f(src1.(string) + src2.(string))
	case func(float64):
		f(src1.(float64) + src2.(float64))
	case func(complex128):
		f(src1.(complex128) + src2.(complex128))
	default:
		return TypeErr{fmt.Errorf("Invalid operation: %v + %v (operator + not defined on %v)", src1, src2, kind)}
	}
	return nil
}
