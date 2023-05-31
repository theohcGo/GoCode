package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	reader := bufio.NewReader(strings.NewReader("http://studygolang.com. \nIt is the home of gophers"))
	line , err := reader.ReadSlice('\n')
	if err != nil {
	}
	// 返回的结果包含界定符本身
	fmt.Printf("the line:%s\n", line)
	n , _ := reader.ReadSlice('\n')
	fmt.Printf("the line:%s\n",line)
	fmt.Printf(string(n));
}
