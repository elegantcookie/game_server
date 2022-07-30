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
	Delete(ctx context.Context, uuid string) error
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
