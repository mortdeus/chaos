package chaos_test

import (
	. "github.com/mortdeus/chaos"
	// "reflect"
	"testing"
)

var InstTests = map[string][]struct {
	src1, src2 interface{}
	dest       interface{}
	errCode    int
}{
	"add": {
		{1, 2, int(0), -1},
		{1.5, 2.3, interface{}(float64(0)), -1},
		{"foo", "bar", interface{}(""), -1},
		{"baz", "qul", "", -1},
		{complex(5.10, 1.532), 5321i, interface{}(complex(1.0, .5)), -1},
		{uint32(512), uint32(1024), uint32(0), -1},

		{uintptr(128), []byte("hi"), interface{}(nil), TypeConv},
		{interface{}(1), interface{}("hi"), interface{}(nil), TypeConv},

		{interface{}("dooble"), interface{}("dobble"), interface{}(""), InvalidOp},
		{true, false, new(bool), InvalidOp},
		{func() { var i = 0; _ = i }, func() { var j = 1; _ = j }, func() {}, InvalidOp},
		{[]int{1, 2, 3}, []int{4, 5, 6}, []int{0}, InvalidOp},
		{map[string]int{"a": 1, "b": 2}, map[string]int{"c": 3, "d": 4}, map[string]int{"": 0}, InvalidOp},
	},
}

func TestAdd(t *testing.T) {
	for _, args := range InstTests["add"] {
		err := Process(Add, args.src1, args.src2, &args.dest)
		if err != nil {
			if e, ok := err.(TypeErr); ok {
				if e.Code() == args.errCode {
					//t.Log(e)
					continue
				}
			}
			t.Error(err)
		}
		t.Logf("%#v + %#v = %#v\n", args.src1, args.src2, args.dest)
	}
}
