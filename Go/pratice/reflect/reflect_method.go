package main

import (
	"log"
	"reflect"
	"strings"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	typ := reflect.TypeOf(&wg)
	//fmt.Println(typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		// NumIn : 方法的参数个数
		argv := make([]string,0,method.Type.NumIn())
		//NumOut : 返回值个数
		returns := make([]string,0,method.Type.NumOut())
		//第0个参数是wg自己
		for j := 1; j < method.Type.NumIn(); j++ {
			// In : 获取方法的第j个输出参数的类型
			// Name : 获取对应类型的名称
			argv = append(argv, method.Type.In(j).Name())
		}

		for j := 0; j < method.Type.NumOut(); j++ {
			returns = append(returns, method.Type.Out(j).Name())
		}

		log.Printf("func (w *%s) %s(%s) %s",typ.Elem().Name(),
												method.Name(),
											strings.Join(argv,","),
											strings.Join(returns,","))
	}
}