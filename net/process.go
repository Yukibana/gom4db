package net

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gom4db/cache"
	"gom4db/pbmessages"
)

func (s *Server)processRequest(request *pbmessages.Request)(responseBuffer []byte){
	switch request.GetType() {
	case pbmessages.REQUEST_MSG_Get_Request:
		return s.processGetRequest(request)
	case pbmessages.REQUEST_MSG_Set_Request:
		return s.processSetRequest(request)
	case pbmessages.REQUEST_MSG_Del_Request:
		return s.processDelRequest(request)
	default:
		return s.UnrecognizedRequestResponse()
	}
}

func  (s *Server) UnrecognizedRequestResponse()[]byte{
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Unknown_Response;
	response.ErrorMsg = "error unsupported request type"
	response.Error = true
	responseBuffer, err := proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func(s *Server)processGetRequest(request *pbmessages.Request)(responseBuffer []byte){
	response := &pbmessages.UnifiedResponse{}
	var key string
	if getRequest := request.GetGetRequest();getRequest != nil{
		key = getRequest.GetKey()
	}else {
		return InvalidFormatResponse(response)
	}
	value, err := s.cache.Get(key)
	sniffError(err)
	if value != nil {
		response.Error = false
		response.Value = cache.Bytes2str(value)
	}else {
		response.Error = true
		response.ErrorMsg = fmt.Sprintf("The value of key {%s} not found",key)
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}


func(s *Server)processSetRequest(request *pbmessages.Request)(responseBuffer []byte){
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Set_Response
	var key,value string
	if setRequest := request.GetSetRequest();setRequest != nil{
		key = setRequest.GetKey()
		value = setRequest.GetValue()
	}else {
		return InvalidFormatResponse(response)
	}
	err := s.cache.Set(key,cache.Str2bytes(value))
	if err != nil {
		response.Error = true
		response.ErrorMsg =  fmt.Sprintf("error occurs when set the key %s ",err.Error())
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}


func(s *Server)processDelRequest(request *pbmessages.Request)(responseBuffer []byte){
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Del_Response
	var key string
	if delRequest := request.GetDelRequest(); delRequest != nil{
		key = delRequest.GetKey()
	}else {
		return InvalidFormatResponse(response)
	}
	err := s.cache.Del(key)
	if err != nil{
		response.Error = true
		response.ErrorMsg =  fmt.Sprintf("error occurs when delete the key %s ",err.Error())
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func InvalidFormatResponse(re *pbmessages.UnifiedResponse)[]byte{
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
