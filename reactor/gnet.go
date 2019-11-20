package reactor

import (
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
}

func (cs *cacheServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumLoops)
	return
}
func (cs *cacheServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	frameData := c.ReadFrame()
	if frameData == nil {
		return
	}
	frameData = append([]byte{}, frameData...)
	_ = cs.workerPool.Submit(func() {
		request := &pbmessages.Request{}
		err := proto.Unmarshal(frameData, request)
		sniffError(err)
		responseBuffer := cs.processRequest(request)
		c.AsyncWrite(responseBuffer)
	})
	return
}

