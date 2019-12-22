package main

import (
	"flag"
	"gom4db/network/cluster"
	"gom4db/network/gonet"
	"gom4db/network/http"
	"gom4db/network/reactor"
)

var gnet bool

func main() {
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.BoolVar(&gnet, "g", false, "Use gnet or not")
	flag.Parse()
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go http.New(n).Listen()
	if gnet {
		reactor.Serve()
	} else {
		gonet.New(n).Listen()
	}
}
