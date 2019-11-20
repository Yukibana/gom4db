package main

import (
	"flag"
	"gom4db/net"
	"gom4db/reactor"
)

var gnet bool

func main() {
	flag.BoolVar(&gnet, "g", false, "Use gnet or not")
	if gnet {
		reactor.Serve()
	} else {
		server := net.New()
		server.Listen()
	}
}
