package main

import (
	"fmt"
	"reflect"
)

// 一个方法的完整信息
type methodType struct {
	method    reflect.Method // 方法的反射对象
	ArgType   reflect.Type   // 方法的入参类型
	ReplyType reflect.Type   // 方法的出参类型
	numCalls  uint64         // 方法被调用的次数
}

type service struct {
	name   string        // 结构体名称
	typ    reflect.Type  // 结构体类型
	rcvr   reflect.Value // 结构体本身 c++中的this指针
	method map[string]*methodType 
}

func NewService(rcvr interface{}) *service {
	s := new(service)

	s.rcvr = reflect.ValueOf(rcvr)
	s.typ = s.rcvr.Type()
	s.name = reflect.Indirect(s.rcvr).Type().Name()

	s.registerMethods()
	return s
}

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		if method.Type.NumIn() != 3 || method.Type.NumOut() != 1 {
			continue
		}
		if method.Type.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType , replyType := mType.In(1) , mType.In(2)
		//TODO：是否是可导出字段

		s.method[method.Name] = &methodType{
			method    : method,
			ArgType   : argType,
			ReplyType : replyType,
		}
	}
}

func (s *service) PrintFun() {
	for k, v := range s.method {
		fmt.Printf("%s %s %s\n",k,v.ArgType,v.ReplyType)
	}
}


type Foo int

type Args struct { Num1 , Num2 int }


func (f Foo) Sum(args Args , reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Pum(args Args , reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}



func main() {
	var f Foo
	s := NewService(&f)
	s.PrintFun()

}