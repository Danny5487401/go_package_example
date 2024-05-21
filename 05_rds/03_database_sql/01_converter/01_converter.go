package main

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

type valueConverterTest struct {
	c   driver.ValueConverter
	in  interface{}
	out interface{}
	err string
}

var now = time.Now()
var answer int64 = 42

type (
	i  int64
	f  float64
	b  bool
	bs []byte
	s  string
	t  time.Time
	is []int
)

var valueConverterTests = []valueConverterTest{
	{driver.Bool, "true", true, ""},
	{driver.Bool, "True", true, ""},
	{driver.Bool, []byte("t"), true, ""},
	{driver.Bool, true, true, ""},
	{driver.Bool, "1", true, ""},
	{driver.Bool, 1, true, ""},
	{driver.Bool, int64(1), true, ""},
	{driver.Bool, uint16(1), true, ""},
	{driver.Bool, "false", false, ""},
	{driver.Bool, false, false, ""},
	{driver.Bool, "0", false, ""},
	{driver.Bool, 0, false, ""},
	{driver.Bool, int64(0), false, ""},
	{driver.Bool, uint16(0), false, ""},
	{c: driver.Bool, in: "foo", err: "sql/driver: couldn't convert \"foo\" into type bool"},
	{c: driver.Bool, in: 2, err: "sql/driver: couldn't convert 2 into type bool"},
	{driver.DefaultParameterConverter, now, now, ""},
	{driver.DefaultParameterConverter, (*int64)(nil), nil, ""},
	{driver.DefaultParameterConverter, &answer, answer, ""},
	{driver.DefaultParameterConverter, &now, now, ""},
	{driver.DefaultParameterConverter, i(9), int64(9), ""},
	{driver.DefaultParameterConverter, f(0.1), float64(0.1), ""},
	{driver.DefaultParameterConverter, b(true), true, ""},
	{driver.DefaultParameterConverter, bs{1}, []byte{1}, ""},
	{driver.DefaultParameterConverter, s("a"), "a", ""},
	{driver.DefaultParameterConverter, is{1}, nil, "unsupported type main.is, a slice of int"},
}

func main() {
	for i, tt := range valueConverterTests {
		out, err := tt.c.ConvertValue(tt.in)
		goterr := ""
		if err != nil {
			goterr = err.Error()
		}
		if goterr != tt.err {
			fmt.Printf("test %d: %T(%T(%v)) error = %q; want error = %q\n",
				i, tt.c, tt.in, tt.in, goterr, tt.err)
		}
		if tt.err != "" {
			continue
		}
		if !reflect.DeepEqual(out, tt.out) {
			fmt.Printf("test %d: %T(%T(%v)) = %v (%T); want %v (%T)\n",
				i, tt.c, tt.in, tt.in, out, out, tt.out, tt.out)
		}
	}
}
