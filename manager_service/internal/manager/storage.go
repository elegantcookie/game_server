package manager

import "context"

type Storage interface {
	Create(ctx context.Context, lr LobbyRecordDTO) (string, error)
	FindById(ctx context.Context, id string) (lr LobbyRecord, err error)
	FindAll(ctx context.Context) (lrs []LobbyRecord, err error)
	Update(ctx context.Context, lr LobbyRecord) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}
