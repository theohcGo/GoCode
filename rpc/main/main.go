package main

import (
	"fmt"
	"geerpc"
	"log"
	"net"
	"sync"
	"time"
)

// 1.定义Foo及Methd
type Foo int

type Args struct { Num1 , Num2 int }

func (f *Foo) Sum(args Args,reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}


type TestMap struct {
	Table []int
	M     map[int]int
}

type ArgsTable struct { Table []int }

func (tm *TestMap) GetMapSize(args ArgsTable,reply *int) error {
	*reply = len(args.Table)
	return nil
}


func startServer(addr chan string) {
	// 2.注册service到server中
	var foo Foo
	if err := geerpc.Register(&foo); err != nil {
		log.Fatal("register foo error: ",err)
	}

	var tm TestMap
	if err := geerpc.Register(&tm); err != nil {
		log.Fatal("register testMap error: ",err)
	}


	l, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:",err)
	}
	log.Println("start rpc server on",l.Addr())
	addr <- l.Addr().String()
	geerpc.Accept(l)
}


func main() {
	log.SetFlags(0)
	addr := make(chan string)
	go startServer(addr)	

	client , err := geerpc.Dial("tcp",<-addr)
	if err != nil {
		fmt.Println("rpc client geerpc.Dial error")
	}
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)
	
	// 同步调用过程
	// send request & receicve response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// 测试foo
			args := &Args{Num1: i,Num2: i*i}
			var reply int
			if err := client.Call("Foo.Sum",args,&reply); err != nil {
				log.Fatal("call Foo.Sum error:",err)
			}
			log.Printf("%d + %d = %d",args.Num1,args.Num2,reply)
			// 测试testMap
			s := []int{1,2,3,45}
			args_tm :=  &ArgsTable{Table: s}
			var reply_tm int
			if err := client.Call("TestMap.GetMapSize",args_tm,&reply_tm); err != nil {
				log.Fatal("call TestMap.GetMapSize error:",err)
			}
			log.Printf("len = %d",reply_tm)
		}(i)
	}
	wg.Wait()
	
	/*
	log.Println("----------------------------------------------------------------------")
	// 异步调用过程
	args := fmt.Sprintf("asyn rpc is hc")
	var reply string
	var reply_c string
	var reply_gwp string
	ch := make(chan *geerpc.Call, 2)
	_ = client.Go("hc.hc", args, &reply, ch)
	_ = client.Go("wp.wp", args, &reply_c, ch)
	_ = client.Go("gwp.gwp", args, &reply_gwp, ch)
	//_ = client.Go("hc.hc", args, &reply, ch)
	//_ = client.Go("hc.hc", args, &reply, ch)
	go func ()  {
		for {
			select {
			case c := <-ch:
				if c.ServiceMethod == "hc.hc" {
					log.Println("reply hc.hc:",reply)
				} else if c.ServiceMethod == "wp.wp" {
					log.Println("reply wp.wp:",reply_c)
				} else if c.ServiceMethod == "gwp.gwp" {
					log.Println("reply gwp.gwp:",reply_gwp)
				}
			}
		}
	}()

	for {
		//runOtherTask()
		time.Sleep(time.Second)
	}
	*/
}

func runOtherTask() {
	log.Println("hello hc")
}