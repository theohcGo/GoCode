package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

// geerpc消息标志值
const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber    int
	CodecType      codec.Type    // 编解码类型
	ConnectTimeout time.Duration // 连接超时
	HandleTimeout    time.Duration // 发送超时
}

// Default Option
var DefaultOption = &Option{
	MagicNumber    : MagicNumber,
	CodecType      :   codec.GobType,
	ConnectTimeout : 10 * time.Second,
	HandleTimeout  : 0,
}

// RPC Server
type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

// Default RPC Server
var DefaultServer = NewServer()

// 接收客户端连接并派发处理协程
func (server *Server) Accept(lis net.Listener) {
	for {
		// test DialTimeout
		//time.Sleep(10 * time.Second)
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		// 开启协程处理连接
		go server.ServeConn(conn)
	}
}

// 协程运行函数
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	var opt Option
	// 解析option字段
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}
	// 判断MagicNumber值是否正确
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number: %x", opt.MagicNumber)
		return
	}
	// 根据CodecType构建Codec实例
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type:%s", opt.CodecType)
		return
	}
	// Codec的处理逻辑
	server.serveCodec(f(conn))
}

// 默认Server开启监听服务
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

// 无效应答
var invalidRequest = struct{}{}

func (server *Server) serveCodec(cc codec.Codec) {
	sendingMux := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		// 读取请求
		req, err := server.readRequest(cc)
		//log.Println("rpc server req :",req.h,req.svc.name,req.mtype.method.Name)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sendingMux)
			continue
		}
		// 处理请求
		// 回复请求
		wg.Add(1)
		go server.handleRequest(cc, req, sendingMux, wg)
	}
	wg.Wait()
	_ = cc.Close()

}

// 请求体封装
type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methodType
	svc          *service
}

// 解析客户端数据并组织为请求头
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		fmt.Println(err)
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// 解析客户端数据并组织为请求体
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	// 生成一个request
	req := &request{h: h}
	// 根据ServiceMethod寻找service对象和methodType对象
	req.svc, req.mtype, err = server.findIndex(h.ServiceMethod)
	if err != nil {
		return req, err
	}

	// req.argv组织请求体的参数
	req.argv = req.mtype.newArgv()
	req.replyv = req.mtype.newReplyv()
	argvi := req.argv.Interface()
	if req.argv.Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	//TODO : argv可能是值类型也可能是指针类型
	//FIX  : ReadBody只能传指针类型
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sendingMux *sync.Mutex) {
	// 保证按顺序发送数据
	sendingMux.Lock()
	defer sendingMux.Unlock()

	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

func (server *Server) handleRequest(cc codec.Codec, req *request, sendingMux *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	// 打印请求头和参数类型
	log.Println(req.h, req.argv)

	callCh := make(chan struct{})
	sendCh := make(chan struct{})
	go func() {
		// 处理请求
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		callCh <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sendingMux)
			sendCh <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sendingMux)
		sendCh <- struct{}{}
	}()
	if timeout == 0 {
		close(callCh)
		close(sendCh)
		return 
	}
	
	select {
	case <-callCh:
		log.Println("rpc server: ")
	case <-sendCh:


	}

}

// 向server中注册一个service
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already define: %s" + s.name)
	}
	return nil
}

// 在server中寻找service.methodType
func (server *Server) findIndex(serviceMethod string) (svc *service, mtype *methodType, err error) {
	//Foo.Sum
	index := strings.LastIndex(serviceMethod, ".")
	if index < 0 {
		err = errors.New("rpc server: service/method request not found" + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:index], serviceMethod[index+1:]
	svcF, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svcF.(*service)
	mtype = svc.method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
		return
	}
	return
}
