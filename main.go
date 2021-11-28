package main

import (
	"goproxy/server"
)

func main() {
	s := server.New()
	s.Run()

}
