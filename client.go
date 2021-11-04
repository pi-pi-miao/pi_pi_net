package pi_pi_net

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

type Client struct {
	netWork string
	Conn    net.Conn
	Addr    string
	Data    chan struct {
		Body []byte
		Err  error
	}
}

func Dail(network, addr string)(*Client, error) {
	c := &Client{
		Data:    make(chan struct{Body []byte
									Err  error},10),
	}
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	c.Conn = conn
	go func() {
		if err := recover(); err != nil {
			log.Println("[client] read data goroutine panic err is :", err)
		}
		c.read()
	}()
	return c,nil
}

func (c *Client) read() {
	size := make([]byte, 2)
	for {
		if _, err := io.ReadFull(c.Conn, size); err != nil {
			c.Data <- struct {
				Body []byte
				Err  error
			}{Err: err}
			continue
		}
		data := make([]byte, binary.LittleEndian.Uint16(size))
		if _, err := io.ReadFull(c.Conn, data); err != nil {
			c.Data <- struct {
				Body []byte
				Err  error
			}{Err: err}
			continue
		}
		c.Data <- struct {
			Body []byte
			Err  error
		}{Body: data,Err: nil}
	}
}

// 客户端读取 返回值为[]byte
// client read and return []byte
func (c *Client) Read() (body []byte,err error){
	data :=<-c.Data
	return data.Body,data.Err
}

// 客户端读取 返回值为string
// client read and return string
func (c *Client) ReadString() (body string,err error){
	data :=<-c.Data
	return ByteToString(data.Body),data.Err
}

// 客户端写入
// client write
func (c *Client) Write(data []byte) (n int, err error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, data...)
	return c.Conn.Write(result)
}

func (c *Client) WriteString(data string) (n int, err error) {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(len(data)))
	result = append(result, StringToByte(data)...)
	return c.Conn.Write(result)
}

func  (c *Client)Close()error{
	return c.Conn.Close()
}
