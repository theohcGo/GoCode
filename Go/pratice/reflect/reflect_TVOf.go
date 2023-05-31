package main

import (
	"fmt"
	"reflect"
)

func main() {
	str := "huchen haha"
	t := reflect.TypeOf(str)
	v := reflect.ValueOf(str)
	fmt.Println("TypeOf str = ",t)
	fmt.Println("VlaueOf str = ",v)
}