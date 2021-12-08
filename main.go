package main

import (
	"goproxy/server"
)

func main() {
	c := server.Options{
		Port: 1080,
		Udp:  true,
	}
	s := server.New(c)
	s.Run()

}
