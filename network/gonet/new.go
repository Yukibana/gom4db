package gonet

import (
	"fmt"
	"gom4db/cache"
	"net"
)

type Server struct {
	cache cache.KeyValueCache
}

func (s *Server) Listen() {
	defer s.cache.Close()
	l, e := net.Listen("tcp", ":12347")
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
func New() *Server {
	fmt.Println("Start Tcp Server without gnet")
	return &Server{cache.NewKeyValueCache()}
}
