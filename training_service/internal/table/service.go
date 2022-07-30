package table

import (
	"context"
	"errors"
	"fmt"
	"training_service/internal/auth"
	"training_service/pkg/logging"
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
	Create(ctx context.Context, dto RecordDTO) (string, error)
	GetAll(ctx context.Context, dto RecordDTO) ([]Record, error)
	GetCollectionNames(ctx context.Context) ([]Collection, error)
	GetById(ctx context.Context, dto RecordDTO) (Record, error)
	GetByUserId(ctx context.Context, dto RecordDTO) (u Record, err error)
	Update(ctx context.Context, dto RecordDTO) error
	Delete(ctx context.Context, dto RecordDTO) error
}

func (s service) Create(ctx context.Context, dto RecordDTO) (recordID string, err error) {
	s.logger.Debug("check password")
	recordID, err = s.storage.Create(ctx, dto)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return recordID, err
		}
		return recordID, fmt.Errorf("failed to create record. error: %w", err)
	}

	return recordID, nil
}

// GetOne record by id
func (s service) GetById(ctx context.Context, dto RecordDTO) (u Record, err error) {
	u, err = s.storage.FindById(ctx, dto)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find record by uuid. error: %w", err)
	}
	return u, nil
}

func (s service) GetByUserId(ctx context.Context, dto RecordDTO) (u Record, err error) {
	u, err = s.storage.FindByUserId(ctx, dto)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return u, err
		}
		return u, fmt.Errorf("failed to find record by userID. error: %w", err)
	}
	return u, nil
}

func (s service) GetAll(ctx context.Context, dto RecordDTO) ([]Record, error) {
	users, err := s.storage.FindAll(ctx, dto)
	if err != nil {
		return users, fmt.Errorf("failed to find records. error: %v", err)
	}
	return users, nil
}

func (s service) GetCollectionNames(ctx context.Context) ([]Collection, error) {
	names, err := s.storage.FindCollectionNames(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get list of collections` names")
	}
	var tableArr []Collection
	for _, name := range names {
		tableArr = append(tableArr, Collection{
			Name: name,
		})
	}
	return tableArr, nil
}

func (s service) Update(ctx context.Context, dto RecordDTO) error {
	err := s.storage.Update(ctx, dto)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update record. error: %w", err)
	}
	return err
}

func (s service) Delete(ctx context.Context, dto RecordDTO) error {
	err := s.storage.Delete(ctx, dto)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete record. error: %w", err)
	}
	return err
}
