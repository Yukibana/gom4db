/*
@Time : 2019/12/21 22:45
@Author : Minus4
*/
package http

import (
	"encoding/json"
	"gom4db/network/cluster"
	"log"
	"net/http"
)

type Server struct {
	cluster.Node
}

func (s *Server)Listen(){
	http.Handle("/cluster",&clusterHandler{s})
	_ = http.ListenAndServe(s.Addr()+":12345", nil)
}

func New(n cluster.Node)*Server{
	return &Server{n}
}

type clusterHandler struct {
	*Server
}

func (h *clusterHandler)ServeHTTP(w http.ResponseWriter,r *http.Request){
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	m := h.Members()
	b, e := json.Marshal(m)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

