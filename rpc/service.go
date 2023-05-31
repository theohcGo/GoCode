package geerpc

import (
	"go/ast"
	"log"
	"reflect"
	"sync/atomic"
)

// 一个方法的完整信息
type methodType struct {
	method     reflect.Method // 方法名
	argType    reflect.Type   // 入参类型
	replyType  reflect.Type   // 出参类型
	numCalls   uint64         // 方法被调用的次数
}

func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}
 
// 创建一个入参实例
func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	if m.argType.Kind() == reflect.Ptr {
		argv = reflect.New(m.argType.Elem())
	} else {
		argv = reflect.New(m.argType).Elem()

	}
	return argv
}

// 创建一个出参实例
func (m *methodType) newReplyv() reflect.Value {
	// replyv一定是指针类型
	replyv := reflect.New(m.replyType.Elem())

	switch m.replyType.Elem().Kind() {
		case reflect.Map:

		case reflect.Slice:

	}
	return replyv
}


type service struct {
	name    string                  // 结构体名称
	typ     reflect.Type            // 结构体类型
	rcvr    reflect.Value
	method  map[string]*methodType  // 结构体所有符合RPC调用条件的方法
}


func newService(rcvr interface{}) *service {
	// 初始化service
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.typ = s.rcvr.Type()
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	
	// 注册该service拥有的方法
	s.registerMethods()
	return s
}	

func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		// 参数个数不符合条件
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		// 返回值不符合条件
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		// 是否是可导出字段
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.method[method.Name] = &methodType{
			method : method,
			argType: argType,
			replyType: replyType,
		}
		log.Printf("rpc server: register %s.%s",s.name,method.Name)
	}
}

// 是否是可导出字段
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}


func (s *service) call(m *methodType, argv,replyv reflect.Value) error {
	atomic.AddUint64(&m.numCalls , 1)
	f := m.method.Func
	returnVals := f.Call([]reflect.Value{s.rcvr,argv,replyv})
	if errInter := returnVals[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}

