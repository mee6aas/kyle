package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RemoveRequested is invoked when the invoker requests to remove an activity.
func (h Handle) RemoveRequested(
	_ context.Context,
	_ string,
	_ string,
) (e error) {
	e = status.Error(codes.Unimplemented, "Out of scope")
	return
}
