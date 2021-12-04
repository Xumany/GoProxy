package server

import (
	"fmt"
	"goproxy/socks"
	"log"
	"net"
	"sync"
)

var Servers = make([]server, 100)

type server struct {
	Connect net.Conn
}

func process(conn net.Conn) {
	if conn == nil {
		return
	}
	// conn.Close()
	var buf [256]byte
	n, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("read from client failed, err:", err)
	}
	socks5 := socks.New(buf[:n])
	b, err := socks5.Auth()
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(b)
	if b[1] == socks.UserPass {
		n, err := conn.Read(buf[:]) // 读取数据 // 将这个数组转换切片进去
		if err != nil {
			fmt.Println("read from client failed, err:", err)
		}
		err = socks.Check(buf[:n])
		if err != nil {
			conn.Write([]byte{buf[0], 1})
			fmt.Println(err)
		}
		conn.Write([]byte{buf[0], 0})
	}
	n, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println(err)
		return
	}

	writer, err := socks.Request(buf[:n])
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	conn.Write(writer.Response)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer writer.Conn.Close()
		socks.Copy(writer.Conn, conn)
	}()

	go func() {
		defer wg.Done()
		defer conn.Close()
		socks.Copy(conn, writer.Conn)
	}()

	wg.Wait()

}
func (s *server) Run() {
	conn, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		log.Print("err")
	}
	defer conn.Close()
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Println("err" + err.Error())
		}
		go process(c)
	}
}

func New() *server {
	return &server{}
}
