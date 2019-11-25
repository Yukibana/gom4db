package gonet

import (
	"gom4db/pbmessages"
	"gom4db/service"
)

func (s *Server) processRequest(request *pbmessages.Request) (responseBuffer []byte) {
	switch request.GetType() {
	case pbmessages.REQUEST_MSG_Get_Request:
		return service.ProcessGetRequest(request, s.cache)
	case pbmessages.REQUEST_MSG_Set_Request:
		return service.ProcessSetRequest(request, s.cache)
	case pbmessages.REQUEST_MSG_Del_Request:
		return service.ProcessDelRequest(request, s.cache)
	default:
		return service.UnrecognizedRequestResponse()
	}
}
