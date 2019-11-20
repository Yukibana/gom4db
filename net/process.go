package net

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gom4db/cache"
	"gom4db/cacheProtoc"
)

func (s *Server)processRequest(request *cacheProtoc.Request)(responseBuffer []byte){
	switch request.GetType() {
	case cacheProtoc.REQUEST_MSG_Get_Request:
		return s.processGetRequest(request)
	case cacheProtoc.REQUEST_MSG_Set_Request:
		return s.processSetRequest(request)
	case cacheProtoc.REQUEST_MSG_Del_Request:
		return s.processDelRequest(request)
	default:
		return s.UnrecognizedRequestResponse()
	}
}

func  (s *Server) UnrecognizedRequestResponse()[]byte{
	response := &cacheProtoc.UnifiedResponse{}
	response.Type = cacheProtoc.RESPONSE_MSG_Unknown_Response;
	response.ErrorMsg = "error unsupported request type"
	response.Error = true
	responseBuffer, err := proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func(s *Server)processGetRequest(request *cacheProtoc.Request)(responseBuffer []byte){
	response := &cacheProtoc.UnifiedResponse{}
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


func(s *Server)processSetRequest(request *cacheProtoc.Request)(responseBuffer []byte){
	response := &cacheProtoc.UnifiedResponse{}
	response.Type = cacheProtoc.RESPONSE_MSG_Set_Response
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


func(s *Server)processDelRequest(request *cacheProtoc.Request)(responseBuffer []byte){
	response := &cacheProtoc.UnifiedResponse{}
	response.Type = cacheProtoc.RESPONSE_MSG_Del_Response
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
