package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
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

func Request(b []byte) ([]byte, error) {
	fmt.Println(b)
	if b[0] != version {
		return nil, errors.New("版本协议不对")
	}
	var atyp = b[3]
	var cmd = b[1]
	var addr string
	var port = binary.BigEndian.Uint16(b[len(b)-2:])
	var c net.Conn
	var requests = []byte{5, 0, 0, 1}
	switch atyp {
	case 1:
		addr = fmt.Sprintf("%d.%d.%d.%d", b[4], b[5], b[6], b[7])
	case 3:
		addr = string(b[5:b[4]])
	default:
		return nil, errors.New("没有该协议")
	}
	var err error
	switch cmd {
	case 1:
		c, err = net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	case 3:
		c, err = net.Dial("udp", fmt.Sprintf("%s:%d", addr, port))
		if err != nil {
			return nil, err
		}
	}
	_ = c.RemoteAddr().String()
	return requests, nil
	// _ := strings.Split(str, ":")
	// var n = str1[0] // IPv4 地址 string
	// ip := strings.Split(n, ".")
	// var cport = str1[1]
	// //var cport = str[1] // 端口号
	// var int16 = make([]uint8, 4)
	// for _, v := range ip {
	// 	vq, err := strconv.Atoi(v)
	// 	vq1 := uint8(vq)
	// 	if err != nil {
	// 		return nil, nil
	// 	}
	// 	int16 = append(int16, vq1)
	// 	// requests = append(requests, byte(v))

	// }
	// cip := int16[4:]
	// newPort, err := strconv.Atoi(cport)
	// if err != nil {
	// 	return nil, err
	// }
	// newPort1 := uint16(newPort)
	// ccc := IntToBytes(newPort1)
	// cip = reverse(cip)
	// ccc = reverse(ccc)
	// requests = append(requests, cip...)
	// requests = append(requests, ccc...)
	// fmt.Println(requests)
	// return requests, nil
}

func IntToBytes(x uint32) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
