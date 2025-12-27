package reporters

import "context"

type Reporter interface {
	Run(ctx context.Context) error
}
