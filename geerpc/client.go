package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/text/cases"
	//"time"
)

// Call : 一次RPC请求的封装
type Call struct {
	Seq           uint64
	ServiceMethod string
	Args          interface{}
	Reply         interface{}
	Error         error
	Done          chan *Call // 异步调用
}

// TODO : 异步通知
func (call *Call) done() {
	call.Done <- call
}

// an RPC Client
type Client struct {
	cc         codec.Codec 
	opt        *Option

    sendingMux sync.Mutex
	mu         sync.Mutex
	header     codec.Header
	seq        uint64
	pending    map[uint64]*Call
	closing    bool
	shutdown   bool
}

// 附带超时机制的Client
type ClientResult struct {
	client *Client
	err    error
}

var _ io.Closer = (*Client)(nil)

var ErrClosing  = errors.New("connection is closing")

// 解析Option字段
func parseOptions(opts ...*Option) (*Option , error) {
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption , nil
	}
	if len(opts) != 1 {
		return nil , errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = DefaultOption.CodecType
	}
	// 设定连接超时时间
	if opt.ConnectTimeout == 0 {
		opt.ConnectTimeout = 10
	}
	return opt , nil
}


func DialTimeout(network , address string, opts ...*Option) (client *Client,err error) {
	// 解析option参数
	opt , err := parseOptions(opts...)
	if err != nil {
		return nil , err
	}
	// TODO : 增加超时
	conn , err := net.DialTimeout(network,address,opt.ConnectTimeout)
	if err != nil {
		log.Println("rpc client: DialTimeout error : ",err)
		return nil , err
	}
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	clientCh := make(chan ClientResult)
	// 开启一个协程创建client
	// 监听Timeout
	go func(clientCh chan ClientResult) {
		select{
		// 超时处理
		case <-time.After(opt.ConnectTimeout):
			log.Println("rpc client : connect to server timeout")
			conn.Close()
		// client成功创建
		case <-clientCh:
			
		}
	}(clientCh)

	select {
	case <-time.After(opt.ConnectTimeout):
		return nil,fmt.Errorf("connectTimeout invalid!")
	case <-clientCh:
		return <-clientCh , nil
	}
	
}


func Dial(network , address string, opts ...*Option) (client *Client,err error) {
	opt , err := parseOptions(opts...)
	if err != nil {
		return nil , err
	}
	conn , err := net.Dial(network,address)
	if err != nil {
		return nil , err
	}
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()
	return NewClient(conn,opt)
}


// 创建一个Client
func NewClient(conn net.Conn,opt *Option) (*Client, error) {
	// 根据CodecType查找Codec的初始化函数
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s",opt.CodecT         ype)
		log.Println("rpc client: codec error:",err)
		return nil,err
	}

	// 和server端协商Option协议
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ",err)
		_ = conn.Close()
		return nil , err
	}
	//time.Sleep(1 * time.Second)
	return newClientCodec(f(conn),opt) , nil
}

// 创建一个ClientCodec编解码器
func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq     : 1   ,
		cc      : cc  ,
		opt     : opt ,
		pending :  make(map[uint64]*Call),
	}
	go client.receive()
	return client
}


// Close client
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing {
		return ErrClosing
	}
	client.closing = true
	return client.cc.Close()
}

// 向pending中注册call
func (client *Client) registerCall(call *Call) (uint64,error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	if client.closing || client.shutdown {
		return 0,ErrClosing
	}
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq , nil
}


// why return *Call
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	//TODO : 判断client是否存在？
	call := client.pending[seq]
	delete(client.pending , seq)
	return call
}

// 服务端或客户端出错，将error通知给所有处于pengding状态的call
func (client *Client) terminateCalls(err error) {
	client.sendingMux.Lock()
	defer client.sendingMux.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

// 接收请求
func (client *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		call := client.removeCall(h.Seq)
		switch {
		// Call不存在
		case call == nil:
			err = client.cc.ReadBody(nil)
		// Call存在但是出错
		case h.Error != "":
			// 传递h.Error
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			// 异步通知本次call调用结束
			call.done()
		// Call存在且正常回应
		default:
			// 解析本次call的响应体
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body" + err.Error())
			}
			// 异步通知本次call调用结束
			call.done()
		}
	}
	client.terminateCalls(err)
}


func (client *Client) send(call *Call) {
	// 确保clien每次send一次完整的请求
	client.sendingMux.Lock()
	defer client.sendingMux.Unlock()

	// register call
	seq, err := client.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return 
	}

	// 准备请求的header
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	if err := client.cc.Write(&client.header,call.Args); err != nil {
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// 异步调用接口
func (client *Client) Go(serviceMethod string , args,reply interface{} , done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call , 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	// 初始化call
	call := &Call {
		ServiceMethod : serviceMethod,
		Args          : args,
		Reply	      : reply,
		Done          : done,
	}
	// TODO OVER : 异步则无需等待send完成
	client.send(call)
	return call
}


// 同步调用接口
func (client *Client) Call(serviceMethod string, args,reply interface{}) error {
	call := <-client.Go(serviceMethod,args,reply,make(chan *Call,1)).Done
	return call.Error
}





