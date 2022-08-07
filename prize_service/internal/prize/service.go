package prize

import (
	"context"
	"errors"
	"fmt"
	"prize_service/internal/auth"
	"prize_service/pkg/logging"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto PrizeDTO) (string, error)
	GetAll(ctx context.Context) ([]Prize, error)
	GetById(ctx context.Context, id string) (Prize, error)
	Update(ctx context.Context, dto Prize) error
	Delete(ctx context.Context, id string) error
}

func (s service) Create(ctx context.Context, dto PrizeDTO) (ticketID string, err error) {
	s.logger.Debug("check password")
	ticket := Prize{
		GameType:  dto.GameType,
		Result:    dto.Result,
		TopPlaces: dto.TopPlaces,
		Reward:    dto.Reward,
		DateTime:  dto.DateTime,
	}
	ticketID, err = s.storage.Create(ctx, ticket)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return ticketID, err
		}
		return ticketID, fmt.Errorf("failed to create lobby. error: %w", err)
	}

	return ticketID, nil
}

// GetOne Prize by id
func (s service) GetById(ctx context.Context, id string) (prize Prize, err error) {
	prize, err = s.storage.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return prize, err
		}
		return prize, fmt.Errorf("failed to find Prize by uuid. error: %w", err)
	}
	return prize, nil
}

func (s service) GetAll(ctx context.Context) ([]Prize, error) {
	prizes, err := s.storage.FindAll(ctx)
	if err != nil {
		return prizes, fmt.Errorf("failed to find prizes. error: %v", err)
	}
	return prizes, nil
}

func (s service) Update(ctx context.Context, prize Prize) error {
	err := s.storage.Update(ctx, prize)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update lobby. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete lobby. error: %w", err)
	}
	return err
}
