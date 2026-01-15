package portfolio

import (
	"context"
	"errors"
	"sync"
)

type memoryRepository struct {
	mu   sync.RWMutex
	data map[string]*Portfolio
}

func NewMemoryRepository(initial []*Portfolio) Repository {
	data := make(map[string]*Portfolio)
	for _, p := range initial {
		data[p.Wallet] = p
	}
	return &memoryRepository{data: data}
}

func (r *memoryRepository) Get(ctx context.Context, wallet string) (*Portfolio, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.data[wallet]
	if !ok {
		return nil, errors.New("portfolio not found")
	}
	return p, nil
}

func (r *memoryRepository) Save(ctx context.Context, p *Portfolio) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[p.Wallet] = p
	return nil
}
