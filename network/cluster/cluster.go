package cluster

import (
	"github.com/hashicorp/memberlist"
	"io/ioutil"
	"stathat.com/c/consistent"
	"time"
)

type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}

type node struct {
	// consistent pointer
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func New(addr, cluster string) (Node, error) {
	conf := memberlist.DefaultLANConfig()
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}

	// clu equals to local addr when clu is empty
	clu := []string{cluster}
	// local join local equals to local
	// local join cluster equals to a bigger cluster
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	circle := consistent.New()
	circle.NumberOfReplicas = 256
	go func() {
		// so it is a kind of polling? fuck you , that's so funny. fuck fuck fuck
		// if one node failed, the others won't get its addr, so the node list will
		// be updated. That's kind of fault
		for {
			// get all members
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			// set consistent circle
			// each node has its own config ?
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}
