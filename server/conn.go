package server

import "net"

func UdpServer() {
	udp, err := net.Listen("udp", ":0")
	if err != nil {
		return
	}
	defer func() {
		udp.Close()
	}()
	for {
		c, err := udp.Accept()
		if err != nil {
			return
		}
		go udpProcess(c)
	}
}
func udpProcess(conn net.Conn) {

}
