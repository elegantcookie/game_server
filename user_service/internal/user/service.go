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
	//Update(ctx context.Context, dto UpdateUserDTO) error
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

// GetOne user by id
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

func (s service) GetAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		return users, fmt.Errorf("failed to find users. error: %v", err)
	}
	return users, nil
}

//func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
//	var updatedUser User
//	s.logger.Debug("compare old and new passwords")
//	if dto.OldPassword != dto.NewPassword || dto.NewPassword == "" {
//		s.logger.Debug("get user by uuid")
//		user, err := s.GetOne(ctx, dto.ID)
//		if err != nil {
//			return err
//		}
//
//		s.logger.Debug("compare hash current password and old password")
//		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.OldPassword))
//		if err != nil {
//			return apperror.BadRequestError("old password does not match current password")
//		}
//
//		dto.Password = dto.NewPassword
//	}
//
//	updatedUser = UpdatedUser(dto)
//
//	s.logger.Debug("generate password hash")
//	err := updatedUser.GeneratePasswordHash()
//	if err != nil {
//		return fmt.Errorf("failed to update user. error %w", err)
//	}
//
//	err = s.storage.Update(ctx, updatedUser)
//
//	if err != nil {
//		if errors.Is(err, apperror.ErrNotFound) {
//			return err
//		}
//		return fmt.Errorf("failed to update user. error: %w", err)
//	}
//	return nil
//}

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
