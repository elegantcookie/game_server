package prize

import "context"

type Storage interface {
	Create(ctx context.Context, prize Prize) (string, error)
	FindById(ctx context.Context, id string) (Prize, error)
	FindAll(ctx context.Context) ([]Prize, error)
	Update(ctx context.Context, prize Prize) error
	Delete(ctx context.Context, id string) error
}
