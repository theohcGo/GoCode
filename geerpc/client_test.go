package geerpc

import (
	"net"
	"strings"
	"testing"
	"time"
)






func TestClient_dialTimeout(t *testing.T) {
	t.Parallel()

	l, _ := net.Listen("tcp", ":0")

	f := func(conn net.Conn, opt *Option) (client *Client, err error) {
		_ = conn.Close()
		time.Sleep(2 * time.Second)
		return nil , nil
	}
 
	// 测试客户端连接超时
	t.Run("timeout",func(t *testing.T) {
		_, err := DialTimeout(f,"tcp",l.Addr().String(),&Option{ ConnectTimeout: time.Second} )
		_assert(err != nil && strings.Contains(err.Error(),"connect timeout"),"expect a timeout error")
	})

	t.Run("0",func(t *testing.T) {
		_, err := DialTimeout(f,"tcp",l.Addr().String(),&Option{ ConnectTimeout: 0} )
		_assert(err != nil, "0 means no limit")
	})


}