package repository

import "context"

type Transaction interface {
	DoInTx(ctx context.Context, f func(ctx context.Context) error) error
}
