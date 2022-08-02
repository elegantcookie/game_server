package ticket

import "context"

type Storage interface {
	Create(ctx context.Context, ticket Ticket) (string, error)
	FindById(ctx context.Context, id string) (Ticket, error)
	FindAll(ctx context.Context) ([]Ticket, error)
	Update(ctx context.Context, ticket Ticket) error
	Delete(ctx context.Context, id string) error
}
