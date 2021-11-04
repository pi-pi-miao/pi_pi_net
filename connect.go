package pi_pi_net

import (
	"encoding/binary"
	"golang.org/x/sys/unix"
	"log"
	"unsafe"
)

type ServerConnect struct {
	Conn  int
	Data  []byte
	Error error
	Ctx   *Context
}

func (ctx *Context) read() {
	defer func() {
		if err := recover();err != nil {
			log.Println("[context] read goroutine panic",err)
		}
	}()
	size := make([]byte, 2)
	for {
		ready, err := unix.EpollWait(ctx.Epoll, ctx.El.Events, ctx.msec)
		if err != nil {
			return
		}
		for i := 0; i < ready; i++ {
			fd := int(ctx.El.Events[i].Fd)
			if _, err := ReadFull(fd, size); err != nil {
				log.Println("read size ", err)
				ctx.Body <- &ServerConnect{
					Conn:  fd,
					Error: err,
					Ctx: ctx,
				}
				continue
			}
			data := make([]byte, binary.LittleEndian.Uint16(size))
			if _, err := ReadFull(fd, data); err != nil {
				log.Println("read file ", err)
				ctx.Body <- &ServerConnect{
					Conn:  fd,
					Error: err,
					Ctx: ctx,
				}
				continue
			}
			ctx.Body <- &ServerConnect{
				Conn: fd,
				Data: data,
				Ctx: ctx,
			}
			continue
		}
	}
}

// 服务端读取
// server read data
func (ctx *Context) ReadServerConnect() *ServerConnect {
	return <-ctx.Body
}

// 阻塞读取
// server block read
func (ctx *Context)ReadBlock()chan *ServerConnect{
	return ctx.Body
}

// 读取数据
// read []byte
func (conn *ServerConnect)ReadByte()[]byte{
	return conn.Data
}

func (conn *ServerConnect)ReadString()string{
	return ByteToString(conn.Data)
}

// 服务端写入 []byte
// server write []byte
func (conn *ServerConnect) Write(data []byte) (n int, err error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, data...)
	return unix.Write(conn.Conn, result)
}

// 服务端写入 string
// server write string
func (conn *ServerConnect) WriteString(data string) (n int, err error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, StringToByte(data)...)
	return unix.Write(conn.Conn, result)
}

// 服务端写入[]byte并且断开连接
// server write []byte and close conn
func (conn *ServerConnect) WriteByteAndClose(data []byte) (writeErr,closeErr error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, data...)
	_,writeErr = unix.Write(conn.Conn, result)
	closeErr = conn.Close()
	return
}

// 服务端写入string并且断开连接
// server write string and close conn
func (conn *ServerConnect) WriteStringAndClose(data string) (writeErr,closeErr error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, StringToByte(data)...)
	_,writeErr = unix.Write(conn.Conn, result)
	closeErr = conn.Close()
	return
}

// 断开连接
// close connect
func (conn *ServerConnect)Close()error{
	m := make(map[int]struct{},len(conn.Ctx.ConnList))
	conn.Ctx.Lock.Lock()
	delete(conn.Ctx.ConnList,conn.Conn)
	for k := range conn.Ctx.ConnList {
		m[k] = struct{}{}
	}
	conn.Ctx.ConnList = m
	conn.Ctx.Lock.Unlock()
	return unix.Close(conn.Conn)
}

func StringToByte(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

func ByteToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

//// string转ytes
//func Str2byte(s string) (b []byte) {
//	*(*string)(unsafe.Pointer(&b)) = s	// 把s的地址付给b
//	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)
//	return
//}
//
//// []byte转string
//func Byte2str(b []byte) string {
//	return *(*string)(unsafe.Pointer(&b))
//}
