package reactor

import (
	"gom4db/pbmessages"
)

func (cs *cacheServer) processRequest(request *pbmessages.Request) (responseBuffer []byte) {
	switch request.GetType() {
	case pbmessages.REQUEST_MSG_Get_Request:
		return processGetRequest(request, cs.cache)
	case pbmessages.REQUEST_MSG_Set_Request:
		return processSetRequest(request, cs.cache)
	case pbmessages.REQUEST_MSG_Del_Request:
		return processDelRequest(request, cs.cache)
	default:
		return unrecognizedRequestResponse()
	}
}
