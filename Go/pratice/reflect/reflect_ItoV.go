package main

import (
	"fmt"
	"reflect"
)

func main() {
	// 接口值 -> 反射值
	str  := "huchen"
	v := reflect.ValueOf(str)
	fmt.Println(v)

	// 反射值 -> 接口值
	a := v.Interface().(string)
	fmt.Println(a)
}