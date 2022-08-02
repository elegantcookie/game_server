package ticket

import (
	"context"
	"errors"
	"fmt"
	"ticket_service/internal/auth"
	"ticket_service/pkg/logging"
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
	Create(ctx context.Context, dto TicketDTO) (string, error)
	GetAll(ctx context.Context) ([]Ticket, error)
	GetById(ctx context.Context, id string) (Ticket, error)
	Update(ctx context.Context, dto Ticket) error
	Delete(ctx context.Context, id string) error
}

func (s service) Create(ctx context.Context, dto TicketDTO) (ticketID string, err error) {
	s.logger.Debug("check password")
	ticket := Ticket{
		Type:         dto.Type,
		TicketPrice:  dto.TicketPrice,
		PlayerAmount: dto.PlayerAmount,
		GameType:     dto.GameType,
		PrizeId:      dto.PrizeId,
	}
	ticketID, err = s.storage.Create(ctx, ticket)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return ticketID, err
		}
		return ticketID, fmt.Errorf("failed to create ticket. error: %w", err)
	}

	return ticketID, nil
}

// GetOne Ticket by id
func (s service) GetById(ctx context.Context, id string) (t Ticket, err error) {
	t, err = s.storage.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return t, err
		}
		return t, fmt.Errorf("failed to find Ticket by uuid. error: %w", err)
	}
	return t, nil
}

func (s service) GetAll(ctx context.Context) ([]Ticket, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to find Tickets. error: %v", err)
	}
	return users, nil
}

func (s service) Update(ctx context.Context, ticket Ticket) error {
	err := s.storage.Update(ctx, ticket)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update Ticket. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete Ticket. error: %w", err)
	}
	return err
}
