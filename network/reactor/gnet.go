package reactor

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool"
	"gom4db/cache"
	"gom4db/pbmessages"
	"log"
)

type cacheServer struct {
	*gnet.EventServer
	addr       string
	multiCore  bool
	async      bool
	codec      gnet.ICodec
	workerPool *pool.WorkerPool
	cache      cache.KeyValueCache
	ch         chan response
	order int
}

func (cs *cacheServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumLoops)
	return
}
func (cs *cacheServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	cs.ch = make(chan response,100)
	err := cs.workerPool.Submit(func() {

	})
	if err != nil {
		panic(err)
	}
	return
}
func (cs *cacheServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	for {
		data := c.ReadFrame()
		if len(data) == 0 {
			return
		}
		// capture iterate
		err := cs.workerPool.Submit(func() {

			frameData := append([]byte{}, data...)
			request := &pbmessages.Request{}
			err := proto.Unmarshal(frameData, request)
			sniffError(err)
			responseBuffer := cs.processRequest(request)
			cs.ch <- response{
				body:  responseBuffer,
				order: cs.order,
			}
			c.AsyncWrite(responseBuffer)
		})
		if err != nil {
			panic(err)
		}
	}
	return
}
func sniffError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
