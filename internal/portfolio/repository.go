package portfolio

import "context"

type Repository interface {
	Get(ctx context.Context, wallet string) (*Portfolio, error)
	Save(ctx context.Context, p *Portfolio) error
}
