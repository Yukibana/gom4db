package reactor

import (
	"container/heap"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
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
	workerPool *goroutine.Pool
	cache      cache.KeyValueCache
}
type cacheContext struct {
	order     int
	responses chan response
}

func (cs *cacheServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("Test codec server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumLoops)
	return
}
func (cs *cacheServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	newContext := cacheContext{
		order:     0,
		responses: make(chan response, 100),
	}
	c.SetContext(&newContext)
	err := cs.workerPool.Submit(func() {
		var responseHeap = new(ResponseHeap)
		heap.Init(responseHeap)
		ctx := c.Context().(*cacheContext)
		currentOrder := 0
		for newResponse := range ctx.responses{
			if newResponse.order == currentOrder {
				c.AsyncWrite(newResponse.body)
				currentOrder++
			} else {
				heap.Push(responseHeap, newResponse)
				continue
			}
			for {
				if responseHeap.IsEmpty() {
					break
				}
				if responseHeap.Top() == currentOrder {
					popResp := heap.Pop(responseHeap).(response)
					c.AsyncWrite(popResp.body)
					currentOrder++
				} else {
					break
				}
			}
		}

	})
	if err != nil {
		panic(err)
	}
	return
}
func (cs *cacheServer)OnClosed(c gnet.Conn,err error)(action gnet.Action){
	ctx := c.Context().(*cacheContext)
	close(ctx.responses)
	return
}
func (cs *cacheServer) React(c gnet.Conn) (out []byte, action gnet.Action) {
	ctx := c.Context().(*cacheContext)
	for {
		data := c.ReadFrame()
		if len(data) == 0 {
			return
		}
		order := ctx.order
		ctx.order++
		err := cs.workerPool.Submit(func() {
			frameData := append([]byte{}, data...)
			request := &pbmessages.Request{}
			err := proto.Unmarshal(frameData, request)
			sniffError(err)
			responseBuffer := cs.processRequest(request)
			ctx.responses <- response{
				body:  responseBuffer,
				order: order,
			}
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
