package invokee

import (
	"context"

	api "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Report reports result for task with given ID.
func Report(ctx context.Context, id string, rst string) (e error) {
	_, e = client.Report(ctx, &api.ReportRequest{
		Id:     id,
		Result: rst,
	})

	return
}
