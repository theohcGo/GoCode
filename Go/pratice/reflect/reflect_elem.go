package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// 非指针类型
	p := Person{Name: "Alice", Age: 30}
	v := reflect.ValueOf(p)
	fmt.Println(v.Kind()) // 输出 struct

	// 指针类型
	p2 := &Person{Name: "Bob", Age: 25}
	v2 := reflect.ValueOf(p2)
	fmt.Println(v2.Kind()) // 输出 ptr

	// 使用Elem()解引用指针类型
	v3 := v2.Elem()
	fmt.Println(v3.Kind()) // 输出 struct

	v4 := reflect.New(reflect.TypeOf(p))
	fmt.Println(v4.Kind())
}