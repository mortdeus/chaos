package chaos

import (
	"fmt"
	"reflect"
)

type Instruction func(src1, src2, dest interface{}) error

func Add(src1, src2, dest interface{}) error {
	sval1, sval2, dval := reflect.ValueOf(src1), reflect.ValueOf(src2), reflect.ValueOf(dest)

	kind := typeCheck(sval1, sval2, dval)
	if dval.Elem().Kind() == reflect.Interface {
		ptr := reflect.New(sval1.Type())
		dval = dval.Elem()
		dval.Set(ptr)

		defer func() {
			dval = reflect.ValueOf(dest)
			i := reflect.Indirect((dval.Elem().Elem().Convert(reflect.PtrTo(sval1.Type())))).Interface()
			dval.Elem().Set(reflect.ValueOf(i))

			//TODO(mortdeus): reenable when verbose mode is implemented
			//fmt.Printf("(%v + %v) = %#v\n", src1, src2, dval.Elem().Interface())
		}()
	}
	switch kind {
	case reflect.Invalid:
		return TypeErr{fmt.Errorf("Incompatible types: \"%v = %v + %v\"", dest, src1, src2)}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		reflect.Indirect(dval).SetInt(sval1.Int() + sval2.Int())

	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		reflect.Indirect(dval.Elem()).SetUint(sval1.Uint() + sval2.Uint())

	case reflect.String:
		reflect.Indirect(dval.Elem()).SetString(sval1.String() + sval2.String())

	case reflect.Float32, reflect.Float64:
		reflect.Indirect(dval.Elem()).SetFloat(sval1.Float() + sval2.Float())

	case reflect.Complex64, reflect.Complex128:
		reflect.Indirect(dval.Elem()).SetComplex(sval1.Complex() + sval2.Complex())

	case reflect.Bool, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Ptr, reflect.Slice, reflect.Struct, reflect.UnsafePointer:

		return TypeErr{fmt.Errorf("Invalid operation: %v + %v (operator + not defined on %v)", src1, src2, kind)}
	}

	return nil
}
