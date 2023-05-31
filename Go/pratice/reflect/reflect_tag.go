package main

import (
	"fmt"
	"reflect"
)

type MyStruct struct {
	Location string `customTag:"custom value"`
}

func main() {
	t := reflect.TypeOf(MyStruct{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag
		fmt.Println(tag)

		customTagValue := tag.Get("customTag")
		fmt.Println(customTagValue)

		lookVal, ok := tag.Lookup("customTag")
		fmt.Println(lookVal," ",ok)
	}

}