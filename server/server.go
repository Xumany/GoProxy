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

//将conn传入到socks 里面验证通过之后 返回所有信息

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
	socks5 := socks.New(buf[:num], conn)
	err = socks5.Auth()
	if err != nil {
		return
	}
	//没有错误就是说明验证通过没有错误
	// 现在需要验证 是什么方法
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

/**
 ** bind 方法
 ** 请求来了之后先创建一个本地端口监听
 ** 将 IP 和端口 使用 socks5 回复命令回复
 ** []byte{5,0,0,IPAaddrtype,IP,端口}
 ** 刚创建的端口来了请求之后
 ** 在将刚刚 来源地址 封装成 socks5命令
 ** []byte{5,0,0,IPAaddrtype,IP,端口}
 ** 回复之后在将数据互相转发
 **/
func BindnMethon() {

}

/**
 ** UDPAssocicte 方法
 ** 传输过来会host和端口
 ** 需要我们绑定的是地址
 ** 都为0服务器自己来控制端口
 ** 链接之后将本地的IP和端口传回客户端
 **/
func UdpAssocicteMethond() {

}

/**
 ** Connect 方法
 **
 **/
func ConnectMethon() {

}

/*
BIND一般用于客户端与服务端建立连接之后，
用于新建立服务端到客户端的连接，
类似于FTP的PORT命令会用到
（有可能记错了，但是肯定是有一个模式是由服务端主动连接客户端，PASV模式？）。
这个过程的标准过程是这样的：
客户端先通过connect建立一个到服务端的信令通道。
客户端再建立一个新的连接到socks5 server，通过bind命令建立一个数据通道。
socks5 server建立一个tcp的监听端口，
但是为了安全起见，bind的时候客户端还会传一个DST地址，
只限于指定地址来临的连接；bind成功之后，
socks5 server返回监听的地址与端口到客户端；
客户端通过信令通道（也可以使用其它任意的通道）把这个监听地址信息传给服务端；
服务端通过正常的tcp connect操作连接socks5 server上的监听端口；
连接建立之后，socks5 server再传送一次BIND成功的信息，
这一次就有了真正的连接者（服务端）的连接信息了。
所以说BIND事件是有两次返回的，
第一次是客户端BIND通过socks5 server建立监听端口成功之后；
第二次是直接连接发生的时候。这段话写得有点绕，主要是在BIND的操作中，
服务端与客户端其实倒了一个身份，变成服务端主动连接客户端了。
UDP associate
UDP用于建立一个UDP的跳转通道（依赖于TCP的socks5协议），过程如下：
客户端与socks5 server建立tcp的连接；
客户端发送udp associate请求，socks5 server分配一个udp的端口，
并将此端口返回给客户端；同前面的BIND一样，
客户端也可以限制使用此udp的使用范围（谁可以使用）；
客户端将自己要发送的UDP报文发送到到socks5 server，
但是前面是有一个复用头的，它告诉socks5 server真正发送到哪里去；
回头的报文也按照相同的格式进行封装。
对于转发的UDP报文，无论是成功与否，都不再有额外的通知报文。
这里的associated的udp通道与tcp通道，具有相同的生命周期。
*/

/*
bind协议 测试工具 FileZilla
Connect/UDP 协议测试工具 QQ SocketTool4
*/
