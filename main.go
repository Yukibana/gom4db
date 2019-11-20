package main

import (
	"flag"
	"gom4db/network/gonet"
	"gom4db/network/reactor"
)

var gnet bool

func main() {
	flag.BoolVar(&gnet, "g", false, "Use gnet or not")
	flag.Parse()
	if gnet{
		reactor.Serve()
	} else {
		gonet.New().Listen()
	}
}
