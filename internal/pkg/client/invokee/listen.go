package invokee

import (
	"context"
	"io"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/mee6aas/zeep/pkg/api/invokee/v1"
)

var (
	onTask chan api.Task
)

func listen(ctx context.Context) (e error) {
	stream, e := client.Listen(ctx, &api.ListenRequest{})
	if e != nil {
		e = errors.Wrap(e, "Failed to request the listen of the invokee service")
		return
	}
	if onTask != nil {
		// warning
		// multiple listen requested
	}

	onTask = make(chan api.Task, 1)

	go func() {
		defer close(onTask)

		for {
			task, e := stream.Recv()
			if e != nil {
				s, _ := status.FromError(e)

				log.WithFields(log.Fields{
					"code":  s.Code(),
					"msg":   s.Message(),
					"err":   s.Err(),
					"isEOF": e == io.EOF,
				}).Debug("Invokee listen stream returns error")

				if s.Code() == codes.PermissionDenied {
					time.Sleep(time.Millisecond * 100)
					stream, _ = client.Listen(ctx, &api.ListenRequest{})
					continue
				}

				return
			}

			onTask <- *task
		}
	}()

	return
}
