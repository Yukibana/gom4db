package main

import (
	"flag"
	"fmt"
	"gom4db/network/cluster"
	"gom4db/network/http"
	"gom4db/network/protobuf"
	"gom4db/network/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
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

	s := grpc.NewServer()
	protobuf.RegisterCacheServiceServer(s, rpc.NewCacheService(n))
	listener, err := net.Listen("tcp", "0.0.0.0:12347")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	reflection.Register(s)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
