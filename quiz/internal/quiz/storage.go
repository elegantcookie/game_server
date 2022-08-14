package quiz

import "context"

type Storage interface {
	Create(ctx context.Context, snake Quiz) (string, error)
	FindById(ctx context.Context, id string) (Quiz, error)
	FindAll(ctx context.Context) ([]Quiz, error)
	Update(ctx context.Context, snake Quiz) error
	Delete(ctx context.Context, id string) error
}
