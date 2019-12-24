package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddRequested is invoked when the invoker requests to add an activity.
func (h Handle) AddRequested(
	_ context.Context,
	_ string,
	_ string,
	_ string,
) (e error) {
	e = status.Error(codes.Unimplemented, "Out of scope")
	return
}
