// +build linux

package pi_pi_net

import (
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net"
	"sync"
)

type EventList struct {
	size   int
	Events []unix.EpollEvent
}

type Context struct {
	Body         chan *ServerConnect
	addr         *net.TCPAddr
	socket       int
	El           *EventList
	Epoll        int
	SetNonblock  bool
	msec         int
	ConnList     map[int]struct{}
	Lock         *sync.RWMutex
}

// 默认使用
func NewContext()*Context{
	return &Context{
		Body:        make(chan *ServerConnect,10),
		El:          newEventList(128),
		SetNonblock: true,
		msec:        -1,
		ConnList: make(map[int]struct{},1000),
		Lock: &sync.RWMutex{},
	}
}

func newEventList(size int) *EventList {
	return &EventList{size, make([]unix.EpollEvent, size)}
}

func (ctx *Context) Listen(network, address string) (*Context, error) {
	socket, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, unix.IPPROTO_TCP)
	if err != nil {
		log.Println("create socket err", err)
		return nil, err
	}
	if err = unix.SetsockoptInt(socket, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
		log.Println("set socket err", err)
		return nil, err
	}
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		log.Println("create address err", err)
		return nil, err
	}
	if err = unix.Bind(socket, &unix.SockaddrInet4{
		Port: tcpAddr.Port,
	}); err != nil {
		log.Println("bind port err", err)
		return nil, err
	}
	if err = unix.Listen(socket, 1<<16-1); err != nil {
		log.Println("listen err", err)
		return nil, err
	}
	ctx.socket = socket
	ctx.addr = tcpAddr
	return ctx, nil
}

func (ctx *Context) Accept() error {
	go ctx.read()
	epoll, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		log.Println("accept err", err)
		return  err
	}
	ctx.Epoll = epoll
	for {
		conn, _, err := unix.Accept(ctx.socket)
		if err != nil {
			log.Println("accept err", err)
			continue
		}
		err = unix.SetNonblock(conn, ctx.SetNonblock)
		if err != nil {
			log.Println("setNonblock err", err)
			continue
		}
		ev := unix.EpollEvent{
			Events: unix.EPOLLPRI | unix.EPOLLIN,
			Fd:     int32(conn),
		}
		if err = unix.EpollCtl(epoll, unix.EPOLL_CTL_ADD, conn, &ev); err != nil {
			log.Println("epoll ctl err", err)
			continue
		}
		ctx.Lock.Lock()
		ctx.ConnList[conn] = struct{}{}
		ctx.Lock.Unlock()
	}
}

func ReadFull(l int, buf []byte) (n int, err error) {
	return readAtLeast(l, buf, len(buf))
}

func readAtLeast(l int, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min {
		var nn int
		nn, err = unix.Read(l, buf[n:])
		if nn < 0 {
			unix.Close(l)
			return
		}
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}
