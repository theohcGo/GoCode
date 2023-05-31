package main

import (
	"fmt"
	"reflect"
)


func Add(a , b int) int {
	return a+b
}

func main() {
	v := reflect.ValueOf(Add)
	if v.Kind() != reflect.Func {
		return 
	}
	t := v.Type() // 根据Value获取Type
	argv := make([]reflect.Value,t.NumIn())
	for i := range argv {
		if t.In(i).Kind() != reflect.Int {
			return 
		}
		argv[i] = reflect.ValueOf(i)
	}
	v2 := v.Call(argv)
	if len(v2) != 1 &&  v2[0].Kind() != reflect.Int {
		return 
	}
	fmt.Println(v2[0].Int())
}