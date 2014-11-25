package chaos_test

import (
	. "github.com/mortdeus/chaos"
	"testing"
)

var InstTests = map[string][]struct {
	src1, src2 interface{}
	dest       interface{}
	errCode    int
}{
	"add": {
		{1, 2, new(int), -1},
		{1.5, 2.3, new(interface{}), -1},
		{"foo", "bar", new(interface{}), -1},
		{"baz", "qul", new(string), -1},
		{complex(5.10, 1.532), 5321i, new(interface{}), -1},
		{uint32(512), uint32(1024), new(uint32), -1},

		{uintptr(128), []byte("hi"), new(interface{}), TypeConv},
		{interface{}(1), interface{}("hi"), new(interface{}), TypeConv},

		{interface{}("dooble"), interface{}("dobble"), new(interface{}), InvalidOp},
		{true, false, new(bool), InvalidOp},
		{func() { var i = 0; _ = i }, func() { var j = 1; _ = j }, new(func()), InvalidOp},
		{true, false, new(bool), InvalidOp},
		{[]int{1, 2, 3}, []int{4, 5, 6}, new([]int), InvalidOp},
		{map[string]int{"a": 1, "b": 2}, map[string]int{"c": 3, "d": 4}, new(map[string]int), InvalidOp},
	},
}

func TestAdd(t *testing.T) {
	for _, args := range InstTests["add"] {
		err := Add(args.src1, args.src2, args.dest)
		if err != nil {
			if e, ok := err.(TypeErr); ok {
				if e.Code() == args.errCode {
					continue
				}
			}
			t.Error(err)
		}
	}
}
