package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, dto UserDTO) error
	FindByUsername(ctx context.Context, username string) (string, error)
}
