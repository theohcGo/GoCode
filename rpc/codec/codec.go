package codec

import "io"
 
// 消息头部
type Header struct {
	ServiceMethod string // RPC调用格式"Service.Method"
	Seq			  uint64 // 消息序列号
	Error 	      string // 错误信息
}

// 消息编解码的接口
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header,interface{}) error
}

// 构造函数
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType Type = "application/gob"
	JsonType Type = "application/json" // 待实现
)

// 不同类型对应不同的构造函数
var NewCodecFuncMap map[Type]NewCodecFunc


func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}





