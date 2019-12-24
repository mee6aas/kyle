package v1

import (
	"context"

	"github.com/mee6aas/zeep/pkg/activity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListRequested is invoked when the invoker requests to list the activities.
func (h Handle) ListRequested(
	_ context.Context,
	_ string,
) (out []activity.Activity, e error) {
	e = status.Error(codes.Unimplemented, "Out of scope")
	return
}
