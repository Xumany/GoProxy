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

type socks struct {
	Version uint8
	Metions uint8
	Metion  []byte
}
type data struct {
	Response []byte
	Conn     net.Conn
}

func New(b []byte) *socks {
	return &socks{
		Version: b[0],
		Metions: b[1],
		Metion:  b[2:],
	}
}

func (s *socks) Auth() ([]byte, error) {
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
	log.Println("Request", b)
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
	var err error
	switch cmd {
	case 1:
		fmt.Println("tcp链接")
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	case 2:
		fmt.Println("bind", b)
		data := &data{
			Response: []byte{5, 0, 0, 1, 127, 0, 0, 1, 16, 22},
		}
		return data, err
	case 3:
		fmt.Println("udp链接")

		if err != nil {
			return nil, err
		}
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

// func DataWiter(dst, src net.Conn) {
// 	var buf = make([]byte, 256)
// 	for {
// 		n, err := src.Read(buf[:])
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 		fmt.Println(buf[:n])
// 		_, err = dst.Write(buf[:n])
// 		if err != nil {
// 			fmt.Println(err)
// 			return
// 		}
// 	}

// }
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
func BindnMethon(port uint16) {

}

/**
 ** UDPAssocicte 方法
 **
 **/
func UdpAssocicteMethond() {

}

/**
 ** Connect 方法
 **
 **/
func ConnectMethon() {

}
