package user

import (
	"context"
	"errors"
	"fmt"
	"user_service/internal/apperror"
	"user_service/pkg/logging"
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
	Create(ctx context.Context, dto CreateUserDTO) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, uuid string) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetByUsernameAndPassword(ctx context.Context, username, password string) (u User, err error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	Delete(ctx context.Context, uuid string) error
	AddTicket(ctx context.Context, dto TicketDTO) error
	DeleteTicket(ctx context.Context, dto TicketDTO) error
}

func (s service) Create(ctx context.Context, dto CreateUserDTO) (userID string, err error) {
	s.logger.Debug("check password")

	user := NewUser(dto)

	s.logger.Debug("generate password hash")
	err = user.GeneratePasswordHash()
	if err != nil {
		s.logger.Errorf("failed to create user due to error %v", err)
		return
	}

	user.Tickets = []string{}
	userID, err = s.storage.Create(ctx, user)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return userID, err
		}
		return userID, fmt.Errorf("failed to create user. error: %w", err)
	}

	return userID, nil
}

// GetById gets user by idd
func (s service) GetById(ctx context.Context, id string) (u User, err error) {
	u, err = s.storage.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find user by uuid. error: %w", err)
	}
	return u, nil
}

func (s service) GetByUsername(ctx context.Context, username string) (u User, err error) {
	u, err = s.storage.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find user by username. error: %w", err)
	}
	return u, nil
}

func (s service) GetByUsernameAndPassword(ctx context.Context, username, password string) (u User, err error) {
	u, err = s.storage.FindByUsername(ctx, username)

	err = u.CheckPassword(password)
	if err != nil {
		return User{}, err
	}
	if err != nil {
		return u, err
	}

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find user by username. error: %w", err)
	}
	return u, nil
}

func (s service) GetAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to find users. error: %v", err)
	}
	return users, nil
}

func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	user := User{
		ID:       dto.ID,
		Username: dto.Username,
		Tickets:  dto.Tickets,
	}
	err := s.storage.Update(ctx, user)
	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, uuid string) error {
	err := s.storage.Delete(ctx, uuid)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return err
}

func ticketInArray(ticketID string, ticketIDS []string) (index int, isIn bool) {
	for i := 0; i < len(ticketIDS); i++ {
		if ticketIDS[i] == ticketID {
			return i, true
		}
	}
	return -1, false
}

func (s service) AddTicket(ctx context.Context, dto TicketDTO) error {
	user, err := s.GetById(ctx, dto.ID)
	if err != nil {
		return err
	}
	if _, isIn := ticketInArray(dto.TicketID, user.Tickets); isIn {
		return fmt.Errorf("user already has ticket with ID: %s", dto.TicketID)
	}

	user.Tickets = append(user.Tickets, dto.TicketID)

	update := UpdateUserDTO{
		ID:       user.ID,
		Username: user.Username,
		Tickets:  user.Tickets,
	}
	err = s.Update(ctx, update)
	if err != nil {
		return err
	}
	return nil
}

func (s service) DeleteTicket(ctx context.Context, dto TicketDTO) error {
	user, err := s.GetById(ctx, dto.ID)
	if err != nil {
		return err
	}
	var (
		index int
		isIn  bool
	)

	if index, isIn = ticketInArray(dto.TicketID, user.Tickets); !isIn {
		return fmt.Errorf("user doesn't have ticket with ID: %s", dto.TicketID)
	}

	user.Tickets = append(user.Tickets[:index], user.Tickets[index+1:]...)

	s.logger.Printf("tickets: %v", user.Tickets)

	update := UpdateUserDTO{
		ID:       user.ID,
		Username: user.Username,
		Tickets:  user.Tickets,
	}
	err = s.Update(ctx, update)
	if err != nil {
		return err
	}
	return nil
}
