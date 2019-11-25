package gonet

import (
	"container/heap"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/smallnest/goframe"
	"gom4db/pbmessages"
	"io"
	"net"
)

func (s *Server) serve(conn net.Conn) {
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
	defer fc.Close()
	responseChan := make(chan response, 1000)
	go s.WriteFrameDaemon(responseChan, fc)

	sequenceNr := 0
	for {
		frameData, err := fc.ReadFrame()
		if err != nil {
			if err == io.EOF {
				return
			} else {
				panic(err)
			}
		}
		go s.AsyncProcessRequest(frameData, sequenceNr, responseChan)
		sequenceNr++
	}
}
func (s *Server) AsyncProcessRequest(frameData []byte, seq int, ch chan response) {
	request := &pbmessages.Request{}
	err := proto.Unmarshal(frameData, request)
	sniffError(err)
	responseBuffer := s.processRequest(request)
	ch <- response{
		body:  responseBuffer,
		order: seq,
	}
}
func (s *Server) WriteFrameDaemon(ch chan response, fc goframe.FrameConn) {
	var responseHeap = new(ResponseHeap)
	heap.Init(responseHeap)

	currentOrder := 0
	for {
		newResponse := <-ch
		if newResponse.order == currentOrder {
			err := fc.WriteFrame(newResponse.body)
			if err != nil {
				sniffError(err)
				return
			}
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
				err := fc.WriteFrame(popResp.body)
				if err != nil {
					sniffError(err)
					return
				}
				currentOrder++
			} else {
				break
			}
		}
	}
}
func sniffError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
