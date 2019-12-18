package main

import (
	"flag"
	"gom4db/network/cluster"
	"gom4db/network/gonet"
	"gom4db/network/reactor"
	"log"
)

var gnet bool

func main() {
	typ := flag.String("type", "inmemory", "cache type")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.BoolVar(&gnet, "g", false, "Use gnet or not")
	log.Println("type is", *typ)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	flag.Parse()
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	if gnet {
		reactor.Serve()
	} else {
		gonet.New(n).Listen()
	}
}
