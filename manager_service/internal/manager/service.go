package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"manager_service/internal/apperror"
	"manager_service/pkg/logging"
	"net/http"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(managerStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: managerStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto LobbyRecordDTO) (string, error)
	GetById(ctx context.Context, id string) (lobby LobbyRecord, err error)
	GetAll(ctx context.Context) ([]LobbyRecord, error)
	UpdateLR(ctx context.Context, lr LobbyRecord) error
	UpdateTime(ctx context.Context, lr LobbyRecord) (int64, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

func (s service) Create(ctx context.Context, dto LobbyRecordDTO) (string, error) {
	lrID, err := s.storage.Create(ctx, dto)
	if err != nil {
		return "", fmt.Errorf("failed to create lr due to: %v", err)
	}
	return lrID, nil
}

func (s service) GetById(ctx context.Context, id string) (lobby LobbyRecord, err error) {
	lobby, err = s.storage.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return lobby, err
		}
		return lobby, fmt.Errorf("failed to find Lobby by uuid. error: %w", err)
	}
	return lobby, nil
}

func (s service) GetAll(ctx context.Context) ([]LobbyRecord, error) {
	lrs, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all lrs due to: %v", err)
	}
	return lrs, nil
}

func (s service) UpdateLR(ctx context.Context, lr LobbyRecord) error {
	log.Println("UPDATE LR")
	err := s.storage.Update(ctx, lr)
	if err != nil {
		return fmt.Errorf("failed to update lr due to: %v", err)
	}
	return nil
}

func (s service) UpdateTime(ctx context.Context, lr LobbyRecord) (int64, error) {
	var u string
	switch lr.Type {
	case lobby:
		u = fmt.Sprintf(updateLobbyTime, lr.LobbyID)
	case qualification:
		u = fmt.Sprintf(updateQualificationTime, lr.GameType)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, u, nil)
	if err != nil {
		return 0, err
	}
	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	if response == nil {
		return 0, fmt.Errorf("response is null")
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var dto UpdateTimeDTO
	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		return 0, err
	}
	return dto.Expiration, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete lr due to: %v", err)
	}
	return nil
}

func (s service) DeleteAll(ctx context.Context) error {
	err := s.storage.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete lr due to: %v", err)
	}
	return nil
}
