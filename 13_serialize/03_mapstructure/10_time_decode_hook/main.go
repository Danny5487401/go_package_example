package main

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}

func Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}

type TestTime struct {
	Name  string
	Birth time.Time // 没有自定义处理函数,会报错: expected a map, got 'string'
}

type Test struct {
	Flag int
	Data interface{}
}

func main() {
	d := TestTime{
		Name:  "Danny",
		Birth: time.Now(),
	}
	test := Test{
		Flag: 1,
		Data: d,
	}
	jsv, _ := json.Marshal(test)
	fmt.Println(string(jsv))
	i := make(map[string]interface{})
	json.Unmarshal(jsv, &i)
	fmt.Println(i)
	r := TestTime{}
	err := Decode(i["Data"], &r)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}
