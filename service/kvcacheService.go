package service

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gom4db/cache"
	"gom4db/pbmessages"
)

func UnrecognizedRequestResponse() []byte {
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Unknown_Response
	response.ErrorMsg = "error unsupported request type"
	response.Error = true
	responseBuffer, err := proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func ProcessGetRequest(request *pbmessages.Request, c cache.KeyValueCache) (responseBuffer []byte) {
	response := &pbmessages.UnifiedResponse{}
	var key string
	if getRequest := request.GetGetRequest(); getRequest != nil {
		key = getRequest.GetKey()
	} else {
		return InvalidFormatResponse(response)
	}
	value, err := c.Get(key)
	sniffError(err)
	if value != nil {
		response.Error = false
		response.Value = cache.Bytes2str(value)
	} else {
		response.Error = true
		response.ErrorMsg = fmt.Sprintf("The value of key {%s} not found", key)
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func ProcessSetRequest(request *pbmessages.Request, c cache.KeyValueCache) (responseBuffer []byte) {
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Set_Response
	var key, value string
	if setRequest := request.GetSetRequest(); setRequest != nil {
		key = setRequest.GetKey()
		value = setRequest.GetValue()
	} else {
		return InvalidFormatResponse(response)
	}
	err := c.Set(key, cache.Str2bytes(value))
	if err != nil {
		response.Error = true
		response.ErrorMsg = fmt.Sprintf("error occurs when set the key %s ", err.Error())
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func ProcessDelRequest(request *pbmessages.Request, c cache.KeyValueCache) (responseBuffer []byte) {
	response := &pbmessages.UnifiedResponse{}
	response.Type = pbmessages.RESPONSE_MSG_Del_Response
	var key string
	if delRequest := request.GetDelRequest(); delRequest != nil {
		key = delRequest.GetKey()
	} else {
		return InvalidFormatResponse(response)
	}
	err := c.Del(key)
	if err != nil {
		response.Error = true
		response.ErrorMsg = fmt.Sprintf("error occurs when delete the key %s ", err.Error())
	}
	responseBuffer, err = proto.Marshal(response)
	sniffError(err)
	return responseBuffer
}

func InvalidFormatResponse(re *pbmessages.UnifiedResponse) []byte {
	re.Error = true
	re.ErrorMsg = "Invalid Format"
	responseBuffer, err := proto.Marshal(re)
	sniffError(err)
	return responseBuffer
}

func sniffError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
