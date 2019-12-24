package main

import (
	"flag"
	"fmt"
	"gom4db/network/cluster"
	"gom4db/network/gonet"
	"gom4db/network/http"
	"gom4db/network/reactor"
	"os"
)

var gnet bool
var node, clus string

func init() {
	flag.StringVar(&node, "node", "", "node address")
	flag.StringVar(&clus, "cluster", "", "cluster address")
	flag.BoolVar(&gnet, "g", false, "Use gnet or not")
	flag.Parse()
	if node == "" {
		if envNode := os.Getenv("NODE"); envNode != "" {
			node = envNode
		} else {
			node = "127.0.0.1"
		}
	}
	if clus == "" {
		if envClus := os.Getenv("CLUSTER"); envClus != "" {
			clus = envClus
		}
	}
}
func main() {
	if clus != "" {
		fmt.Println("Cluster is ", clus)
	}
	fmt.Println("Node is ", node)
	n, e := cluster.New(node, clus)
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
