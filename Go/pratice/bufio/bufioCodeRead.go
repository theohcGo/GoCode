type Reader struct {
	buf          []byte        // 缓存
	rd           io.Reader    // 底层的io.Reader
	// r:从buf中读走的字节（偏移）；w:buf中填充内容的偏移；
	// w - r 是buf中可被读的长度（缓存数据的大小），也是Buffered()方法的返回值
	r, w         int
	err          error        // 读过程中遇到的错误
	lastByte     int        // 最后一次读到的字节（ReadByte/UnreadByte)
	lastRuneSize int        // 最后一次读到的Rune的大小 (ReadRune/UnreadRune)
}


func NewReader(rd io.Reader) *Reader {
	// 默认缓存大小：defaultBufSize=4096  4K
	return NewReaderSize(rd, defaultBufSize)
}


func NewReaderSize(rd io.Reader, size int) *Reader {
	// 已经是bufio.Reader类型，且缓存大小不小于 size，则直接返回
	// 接口类型断言
	b, ok := rd.(*Reader)
	if ok && len(b.buf) >= size {
		return b
	}
	// 缓存大小不会小于 minReadBufferSize （16字节）
	if size < minReadBufferSize {
		size = minReadBufferSize
	}
	// 构造一个bufio.Reader实例
	return &Reader{
		buf:          make([]byte, size),
		rd:           rd,
		lastByte:     -1,
		lastRuneSize: -1,
	}
}