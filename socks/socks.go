package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
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
	data    []byte
}
type data struct {
	Response []byte
	Conn     net.Conn
}

func New(b []byte) *socks {

	return &socks{
		Version: b[0],
		Metions: b[1],
		data:    b[2:],
	}
}

func (s *socks) Auth() ([]byte, error) {
	fmt.Println(s.data[2 : s.Metions+2])
	if s.Version != version {
		return nil, errors.New("协议版本错误")
	}
	for _, v := range s.data {
		if v == NoUser {
			return []byte{version, NoUser}, nil
		}
		if v == UserPass {
			fmt.Println("账户密码认证")
			return []byte{version, UserPass}, nil
		}
	}

	return nil, errors.New("协议不对")
}
func Check(b []byte) error {
	var u = string(b[2 : b[1]+1])
	var p = string(b[2+b[1]:])
	fmt.Print("账号:" + u + "\n" + "密码:" + p)
	if u == username && p == password {
		return nil
	}
	return errors.New("账号密码错误")
}

func Request(b []byte) (*data, error) {
	fmt.Println(b)
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
		fmt.Printf("DNS解析前%s 解析后的地址%s\n", dmian, addr)
	default:
		return nil, errors.New("没有该协议")
	}
	var err error
	switch cmd {
	case 1:
		conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}

	case 3:
		conn, err = net.Dial("udp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	}
	str := conn.RemoteAddr().String()
	str1 := strings.Split(str, ":")
	str = str1[0]
	Port2 := str1[1]
	str1 = strings.Split(str, ".")
	for _, v := range str1 {
		n, _ := strconv.Atoi(v)
		x := uint8(n)
		requests = append(requests, x)

	}
	num, _ := strconv.Atoi(Port2)
	c := uint16(num)
	bit := uin16ToBigendBytes(c)

	requests = append(requests, bit...)
	fmt.Printf("requests: %v\n", requests)
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
