package net

import (
	"gom4db/pbmessages"
)
func (s *Server) processRequest(request *pbmessages.Request) (responseBuffer []byte) {
	switch request.GetType() {
	case pbmessages.REQUEST_MSG_Get_Request:
		return processGetRequest(request, s.cache)
	case pbmessages.REQUEST_MSG_Set_Request:
		return processSetRequest(request, s.cache)
	case pbmessages.REQUEST_MSG_Del_Request:
		return processDelRequest(request, s.cache)
	default:
		return unrecognizedRequestResponse()
	}
}
