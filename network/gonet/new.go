package gonet

import (
	"fmt"
	"gom4db/cache"
	"gom4db/network/cluster"
	"net"
)

type Server struct {
	cache cache.KeyValueCache
	cluster.Node
}

func (s *Server) Listen() {
	defer s.cache.Close()
	fmt.Println("Listening at tcp "+s.Addr()+":12347")
	l, e := net.Listen("tcp", s.Addr()+":12347")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.serve(c)
	}
}
func New(n cluster.Node) *Server {
	fmt.Println("start tcp server in gonet mode")
	return &Server{cache.NewKeyValueCache(),n}
}
