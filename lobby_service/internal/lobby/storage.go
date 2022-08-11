package lobby

import "context"

type Storage interface {
	Create(ctx context.Context, lobby Lobby) (string, error)
	FindById(ctx context.Context, id string) (Lobby, error)
	FindByParams(ctx context.Context, gameType string, maxPlayers, prizeSum int) (string, error)
	FindAll(ctx context.Context) ([]Lobby, error)
	Update(ctx context.Context, lobby Lobby) error
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}
