package socks

import (
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
			fmt.Println("无需认证")
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

func Request(c net.IP) {
	// net.DialTCP("tcp4", nil, )
}
