package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	version  = uint8(5)
	NoUser   = uint8(0)
	UserPass = uint8(2)
	username = "admin"
	password = "admin"
)

type socksAuth struct {
	Version uint8
	Metions uint8
	Metion  []byte
}

type data struct {
	Response []byte
	Conn     net.Conn
}

func New(b []byte) *socksAuth {
	return &socksAuth{
		Version: b[0],
		Metions: b[1],
		Metion:  b[2:],
	}
}

func (s *socksAuth) Auth() ([]byte, error) {
	if s.Version != version {
		return nil, errors.New("协议版本错误")
	}
	for _, v := range s.Metion {
		if v == NoUser {
			return []byte{version, NoUser}, nil
		}
		if v == UserPass {
			return []byte{version, UserPass}, nil
		}
	}
	return nil, errors.New("协议不对")
}

// 验证账号密码
func Check(b []byte) error {
	var u = string(b[2 : b[1]+1])
	var p = string(b[2+b[1]:])

	if u == username && p == password {
		return nil
	}
	return errors.New("账号密码错误")
}

func Request(b []byte) (*data, error) {
	if b[0] != version {
		return nil, errors.New("版本协议不对")
	}
	var atyp = b[3]
	var cmd = b[1]
	var addr string
	var port = binary.BigEndian.Uint16(b[len(b)-2:])
	var conn net.Conn
	var requests = []byte{5, 0, 0, 1}
	switch atyp {
	case 1:
		addr = fmt.Sprintf("%d.%d.%d.%d", b[4], b[5], b[6], b[7])
	case 3:
		dmian := string(b[5 : len(b)-2])
		ip, _ := net.ResolveIPAddr("ip", dmian)
		addr = ip.String()
		fmt.Printf("DNS解析前%s 解析后的地址%s\n 端口:%d", dmian, addr, port)
	default:
		return nil, errors.New("没有该协议")
	}
	fmt.Println(addr)
	var err error
	switch cmd {
	case 1:
		ConnectMethon()
		fmt.Println("tcp链接")
	case 2:
		BindnMethon()
		fmt.Println("bind")
		return nil, err
	case 3:
		UdpAssocicteMethond()
		fmt.Println("udp链接")
	}
	str := conn.RemoteAddr().String()

	bit, _ := iPToByte(str)
	requests = append(requests, bit...)
	log.Println("Response:", requests)
	data := data{
		Response: requests,
		Conn:     conn,
	}
	return &data, nil
}

func uin16ToBigendBytes(num uint16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, num)
	return bytesBuffer.Bytes()
}

func Copy(src io.ReadWriteCloser, dst io.ReadWriteCloser) (written int64, err error) {
	size := 1024
	buf := make([]byte, size)
	for {
		nr, e := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}

		if e != nil {
			if e != io.EOF {
				err = e
			}
			break
		}
	}
	return written, err
}
func iPToByte(str string) (b []byte, err error) {
	ip, port, err := net.SplitHostPort(str)
	if err != nil {
		return nil, err
	}
	str1 := strings.Split(ip, ".")
	p, err := strconv.Atoi(port)
	port1 := uint16(p)
	portS := uin16ToBigendBytes(port1)
	for _, v := range str1 {
		n, _ := strconv.Atoi(v)
		x := uint8(n)
		b = append(b, x)
	}
	b = append(b, portS...)
	return b, err
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
