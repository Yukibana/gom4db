package reactor

import (
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool"
	"gom4db/cache"
	"time"
)

func Serve() {
	var port int
	var multiCore bool
	// Example command: go run server.go --port 9000 --multiCore true
	flag.IntVar(&port, "port", 12347, "server port")
	flag.BoolVar(&multiCore, "multiCore", true, "multiCore")
	flag.Parse()
	addr := fmt.Sprintf("tcp://:%d", port)

	encoderConfig := gnet.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	decoderConfig := gnet.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}
	codec := gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)
	cs := &cacheServer{addr: addr, multiCore: multiCore, async: true, codec: codec, workerPool: pool.NewWorkerPool(),cache:cache.NewKeyValueCache()}

	err := gnet.Serve(cs, addr, gnet.WithMulticore(multiCore), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(codec))
	if err != nil {
		panic(err)
	}
}
