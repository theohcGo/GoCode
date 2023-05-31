package main

import (
	"fmt"
	"reflect"
)

func main() {
	s := "strGo"
	s_fv := reflect.ValueOf(&s)
	fmt.Println(s)
	s_fv.Elem().SetString("strHC")
	fmt.Println(s)
}