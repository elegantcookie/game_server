package snake

import (
	"context"
	"errors"
	"fmt"
	"snake_service/internal/auth"
	"snake_service/pkg/logging"
	"time"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(storage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto SnakeDTO) (string, error)
	GetAll(ctx context.Context) ([]Snake, error)
	GetById(ctx context.Context, id string) (Snake, error)
	Update(ctx context.Context, dto Snake) error
	Delete(ctx context.Context, id string) error
	SendResult(ctx context.Context, dto SendResultDTO) error
	GetGameStatus(ctx context.Context, gsID string) (int, error)
}

func (s service) Create(ctx context.Context, dto SnakeDTO) (snakeID string, err error) {
	snake := NewSnake(dto)
	snakeID, err = s.storage.Create(ctx, snake)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return snakeID, err
		}
		return snakeID, fmt.Errorf("failed to create game server. error: %w", err)
	}

	return snakeID, nil
}

// GetById get game server data by id
func (s service) GetById(ctx context.Context, id string) (t Snake, err error) {
	t, err = s.storage.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return t, err
		}
		return t, fmt.Errorf("failed to find game server by uuid. error: %w", err)
	}
	return t, nil
}

func (s service) GetAll(ctx context.Context) ([]Snake, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to find game servers. error: %v", err)
	}
	return users, nil
}

func (s service) Update(ctx context.Context, snake Snake) error {
	err := s.storage.Update(ctx, snake)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update Snake. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete Snake. error: %w", err)
	}
	return err
}

func (s service) SendResult(ctx context.Context, dto SendResultDTO) error {
	gs, err := s.storage.FindById(ctx, dto.GameServerID)
	if err != nil {
		return fmt.Errorf("failed to find game server due to: %v", err)
	}
	var isIn bool
	for i, player := range gs.Results {
		if player.UserID == dto.UserID {
			isIn = true
			gs.Results[i].Result = dto.Result
		}
	}
	if !isIn {
		p := NewPlayer(dto)
		gs.Results = append(gs.Results, p)
	}
	err = s.Update(ctx, gs)
	if err != nil {
		return fmt.Errorf("failed to update game server due to: %v", err)
	}
	return nil
}

// GetGameStatus returns game status.
// Status values and their description are in consts.go file.
func (s service) GetGameStatus(ctx context.Context, gsID string) (int, error) {
	gs, err := s.GetById(ctx, gsID)
	if err != nil {
		return StatusError, err
	}
	nowTimestamp := time.Now().Unix()
	if nowTimestamp >= gs.EndTime {
		return StatusEnded, nil
	}
	if nowTimestamp >= gs.StartTime {
		return StatusStarted, nil
	}
	return StatusNotStarted, nil
}
