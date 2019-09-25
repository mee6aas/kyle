package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/mee6aas/zeep/pkg/api/invokee/v1"

	assigns "github.com/mee6aas/kyle/internal/pkg/var/assignments"
)

// Reported is invoked when the invokee client requested to report for the task.
func (h Handle) Reported(req *v1.ReportRequest) (e error) {
	id := req.GetId()

	if ok := assigns.Report(id, req); !ok {
		e = status.Error(codes.NotFound, "Invocation ID not found")
		return
	}

	return
}
