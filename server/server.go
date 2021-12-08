package server

import (
	"fmt"
	"goproxy/socks"
	"log"
	"net"
)

type Options struct {
	Port     uint16
	Udp      bool
	Bind     bool
	User     string
	Pass     string
	LogLevel string
}

func process(conn net.Conn) {
	if conn == nil {
		return
	}
	// conn.Close()
	var buf [50]byte
	num, err := conn.Read(buf[:])
	if err != nil {
		log.Println("read from client failed, err:", err)
	}
	socks5 := socks.New(buf[:num])
	b, err := socks5.Auth()
	if err != nil {
		return
	}
	fmt.Println(b)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// conn.Write(b)
	// if b[1] == socks.UserPass {
	// 	n, err := conn.Read(buf[:]) // 读取数据 // 将这个数组转换切片进去
	// 	if err != nil {
	// 		fmt.Println("read from client failed, err:", err)
	// 	}
	// 	err = socks.Check(buf[:n])
	// 	if err != nil {
	// 		conn.Write([]byte{buf[0], 1})
	// 		fmt.Println(err)
	// 	}
	// 	conn.Write([]byte{buf[0], 0})
	// }
	// n, err = conn.Read(buf[:])
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// writer, err := socks.Request(buf[:n])
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// conn.Write(writer.Response)

	//	go func() {
	//		defer writer.Conn.Close()
	//		socks.Copy(writer.Conn, conn)
	//	}()

	//	go func() {
	//		defer conn.Close()
	//		socks.Copy(conn, writer.Conn)
	//	}()

}
func (s *Options) Run() error {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Println("err:", err)
	}
	defer conn.Close()
	for {
		c, err := conn.Accept()
		if err != nil {
			log.Println("err:", err)
		}
		go process(c)
	}
}

func New(o Options) *Options {
	return &Options{}
}
