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

	user.Tickets = []GameTickets{}
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

func getGameTickets(gameType string, gameTicketsArr []GameTickets) (index int, found bool) {
	for i := 0; i < len(gameTicketsArr); i++ {
		if gameTicketsArr[i].GameType == gameType {
			return i, true
		}
	}
	return -1, false
}

func getTicketID(ticketID string, gameTickets GameTickets) (index int, found bool) {
	for i := 0; i < len(gameTickets.IDsOfGT); i++ {
		if gameTickets.IDsOfGT[i] == ticketID {
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

	i, gtFound := getGameTickets(dto.GameType, user.Tickets)
	if gtFound {
		_, found := getTicketID(dto.TicketID, user.Tickets[i])
		if found {
			return fmt.Errorf("ticket is already in list")
		}
		user.Tickets[i].IDsOfGT = append(user.Tickets[i].IDsOfGT, dto.TicketID)
		user.Tickets[i].Amount += 1
	} else {
		user.Tickets = append(user.Tickets, GameTickets{
			GameType: dto.GameType,
			Amount:   1,
			IDsOfGT:  []string{dto.TicketID},
		})
	}

	update := UpdateUserDTO{
		ID:            user.ID,
		Username:      user.Username,
		HasFreeTicket: user.HasFreeTicket,
		Tickets:       user.Tickets,
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

	i, gtFound := getGameTickets(dto.GameType, user.Tickets)
	if gtFound {
		j, found := getTicketID(dto.TicketID, user.Tickets[i])
		if !found {
			return fmt.Errorf("ticket with id: %s not found", dto.TicketID)
		}
		user.Tickets[i].IDsOfGT = append(user.Tickets[i].IDsOfGT[:j], user.Tickets[i].IDsOfGT[j+1:]...)
		user.Tickets[i].Amount -= 1
	} else {
		return fmt.Errorf("user has no tickets of game type: %s", dto.GameType)
	}

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
