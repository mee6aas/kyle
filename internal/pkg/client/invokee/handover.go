package invokee

import (
	"context"

	api "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

// Handover notifies that this invokee will handover the its control.
func Handover(ctx context.Context) (e error) {
	_, e = client.Handover(ctx, &api.HandoverRequest{})

	return
}
