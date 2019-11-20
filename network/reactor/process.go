package reactor

import (
	"gom4db/pbmessages"
	"gom4db/service"
)

func (cs *cacheServer) processRequest(request *pbmessages.Request) (responseBuffer []byte) {
	switch request.GetType() {
	case pbmessages.REQUEST_MSG_Get_Request:
		return service.ProcessGetRequest(request, cs.cache)
	case pbmessages.REQUEST_MSG_Set_Request:
		return service.ProcessSetRequest(request, cs.cache)
	case pbmessages.REQUEST_MSG_Del_Request:
		return service.ProcessDelRequest(request, cs.cache)
	default:
		return service.UnrecognizedRequestResponse()
	}
}
