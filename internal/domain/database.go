package domain

import "context"

type DatabaseTransaction interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
