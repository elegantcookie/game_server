package snake

import "context"

type Storage interface {
	Create(ctx context.Context, snake Snake) (string, error)
	FindById(ctx context.Context, id string) (Snake, error)
	FindAll(ctx context.Context) ([]Snake, error)
	Update(ctx context.Context, snake Snake) error
	Delete(ctx context.Context, id string) error
}
