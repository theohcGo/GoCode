package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser
	buf *bufio.Writer // 应答缓冲区
	dec *gob.Decoder  // 解码对象
	enc *gob.Encoder  // 编码对象
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn : conn,
		buf  : buf,
		dec  : gob.NewDecoder(conn),
		enc  : gob.NewEncoder(buf),
	}
}

// 实现Codec接口
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// body为传出参数
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header,body interface{}) (err error) {
	// 确保Write函数返回前刷新缓冲区并关闭连接	
	defer func() {
		// 数据通过网络发送
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	// 编码的数据放在buf中
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header" , err)
		return 
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body" , err)
		return 
	}
	return 
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
