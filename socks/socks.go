package socks

import (
	"bytes"
	"encoding/binary"
	"errors"
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
	Conn    net.Conn
	Version uint8
	Metions uint8
	Metion  []byte
}
type socksUser struct {
	User string
	Pass string
}

func New(b []byte, Conn net.Conn) *socks {
	return &socks{
		Conn:    Conn,
		Version: b[0],
		Metions: b[1],
		Metion:  b[2:],
	}
}

func (s *socks) Auth() error {
	if s.Version != version {
		s.Conn.Write([]byte{5, 0xff})
		return errors.New("version Error ")
	}
	for _, v := range s.Metion {
		if v == NoUser {
			s.Conn.Write([]byte{version, NoUser})
			return nil
		}
		if v == UserPass {
			s.Conn.Write([]byte{version, UserPass})
			c, err := getUserPass(s.Conn)
			if err != nil {
				return err
			}
			err = c.Check()
			if err != nil {
				return err
			}
			return nil
		}
	}
	s.Conn.Write([]byte{5, 0xff})
	return errors.New("协议不对")
}
func getUserPass(c net.Conn) (socksUser, error) {
	var buff = make([]byte, 50)
	n, err := c.Read(buff)
	if err != nil {
		return socksUser{}, err
	}
	buff = buff[:n]
	s := socksUser{
		User: string(buff[2 : buff[1]+1]),
		Pass: string(buff[2+buff[1]:]),
	}
	return s, nil
}

// 验证账号密码
func (s *socksUser) Check() error {

	if s.User == username && s.Pass == password {
		return nil
	}
	return errors.New("账号密码错误")
}

func uin16ToBigendBytes(num uint16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, num)
	return bytesBuffer.Bytes()
}

// func Copy(src io.ReadWriteCloser, dst io.ReadWriteCloser) (written int64, err error) {

// 	buf := make([]byte, 1024)
// 	for {
// 		nr, e := src.Read(buf)
// 		if nr > 0 {
// 			nw, ew := dst.Write(buf[0:nr])
// 			if nw > 0 {
// 				written += int64(nw)
// 			}
// 			if ew != nil {
// 				err = ew
// 				break
// 			}
// 			if nr != nw {
// 				err = io.ErrShortWrite
// 				break
// 			}
// 		}

// 		if e != nil {
// 			if e != io.EOF {
// 				err = e
// 			}
// 			break
// 		}
// 	}
// 	return written, err
// }
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
