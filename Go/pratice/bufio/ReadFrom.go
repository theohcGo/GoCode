package main

import (
	"io"
	"os"
	"fmt"
)

func ReadFrom(reader io.Reader , num int) ([]byte , error) {
	p := make([]byte , num)
	n , err := reader.Read(p)
	if n > 0 {
		return p[:n] , nil
	}
	return p , err
}


func main() {
	//从标准输入中读取
	data , err := ReadFrom(os.Stdin , 11)
	if err != nil {
		fmt.Printf("%T\n",err);
	}
	fmt.Print(data)
}