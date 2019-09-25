package v1

import (
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HandoverRequested not supported for this agent delegate.
func (h Handle) HandoverRequested(addr *net.TCPAddr) (e error) {
	e = status.Error(codes.Unimplemented, "Unimplemented")

	return
}
