package net

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/smallnest/goframe"
	"gom4db/cacheProtoc"
	"io"
	"net"
)
func (s *Server) serve(conn net.Conn) {
	defer conn.Close()
	encoderConfig := goframe.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}

	decoderConfig := goframe.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 4,
	}
	fc := goframe.NewLengthFieldBasedFrameConn(encoderConfig, decoderConfig, conn)

	for{
		frameData,err := fc.ReadFrame()
		if err != nil{
			if err == io.EOF{
				return
			}else {
				panic(err)
			}
		}
		request := &cacheProtoc.Request{}
		err = proto.Unmarshal(frameData, request)
		sniffError(err)
		responseBuffer := s.processRequest(request)
		err = fc.WriteFrame(responseBuffer)
		if err != nil{
			sniffError(err)
			return
		}
	}
}

func InvalidFormatResponse(re *cacheProtoc.UnifiedResponse)[]byte{
	re.Error = true
	re.ErrorMsg = "Invalid Format"
	responseBuffer, err := proto.Marshal(re)
	sniffError(err)
	return responseBuffer
}

func sniffError(err error) {
	if err != nil{
		fmt.Println(err)
	}
}
