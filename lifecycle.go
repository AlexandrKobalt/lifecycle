package lifecycle

import "context"

type Lifecycle interface {
	Start(ctx context.Context) (err error)
	Stop(ctx context.Context) (err error)
}
