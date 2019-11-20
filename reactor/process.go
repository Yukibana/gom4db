package reactor

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gom4db/cache"
	"gom4db/cacheProtoc"
)

func (cs *cacheServer)processRequest(request *cacheProtoc.Request)(responseBuffer []byte){
	response := &cacheProtoc.UnifiedResponse{}
	switch request.GetType() {
	case cacheProtoc.REQUEST_MSG_Get_Request:
		response.Type = cacheProtoc.RESPONSE_MSG_Get_Response
		var key string
		if getRequest := request.GetGetRequest();getRequest != nil{
			key = getRequest.GetKey()
		}else {
			return InvalidFormatResponse(response)
		}
		value, err := cs.cache.Get(key)
		sniffError(err)
		if value != nil {
			response.Error = false
			response.Value = cache.Bytes2str(value)
		}else {
			response.Error = true
			response.ErrorMsg = fmt.Sprintf("The value of key {%s} not found",key)
		}
	case cacheProtoc.REQUEST_MSG_Set_Request:
		response.Type = cacheProtoc.RESPONSE_MSG_Set_Response
		var key,value string
		if setRequest := request.GetSetRequest();setRequest != nil{
			key = setRequest.GetKey()
			value = setRequest.GetValue()
		}else {
			return InvalidFormatResponse(response)
		}
		err := cs.cache.Set(key,cache.Str2bytes(value))
		if err != nil {
			response.Error = true
			response.ErrorMsg =  fmt.Sprintf("error occurs when set the key %s ",err.Error())
		}
	case cacheProtoc.REQUEST_MSG_Del_Request:
		response.Type = cacheProtoc.RESPONSE_MSG_Del_Response
		var key string
		if delRequest := request.GetDelRequest(); delRequest != nil{
			key = delRequest.GetKey()
		}else {
			return InvalidFormatResponse(response)
		}
		err := cs.cache.Del(key)
		if err != nil{
			response.Error = true
			response.ErrorMsg =  fmt.Sprintf("error occurs when delete the key %s ",err.Error())
		}
	default:
		response.Type = cacheProtoc.RESPONSE_MSG_Unknown_Response;
		response.ErrorMsg = "error unsupported request type"
		response.Error = true
	}
	responseBuffer, err := proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}
func InvalidFormatResponse(re *cacheProtoc.UnifiedResponse)[]byte{
	re.Error = true
	re.ErrorMsg = "Invalid Format"
	responseBuffer, err := proto.Marshal(re)
	sniffError(err)
	return responseBuffer
}