package table

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"training_service/internal/auth"
	"training_service/internal/config"
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
	CreateCollection(ctx context.Context, dto CollectionDTO) error
	DeleteCollection(ctx context.Context, dto CollectionDTO) error
	GetAll(ctx context.Context, dto RecordDTO) ([]Record, error)
	GetCollectionNames(ctx context.Context) ([]Collection, error)
	GetById(ctx context.Context, dto RecordDTO) (Record, error)
	GetByUserId(ctx context.Context, dto RecordDTO) (u Record, err error)
	Delete(ctx context.Context, dto RecordDTO) error
	Update(ctx context.Context, dto RecordDTO) error
	UpdateTime(ctx context.Context, tableName string) (int64, error)
}

func (s service) CreateCollection(ctx context.Context, dto CollectionDTO) error {
	s.logger.Debug("create collection")

	err := s.storage.CreateCollection(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to create record. error: %w", err)
	}
	err = s.NotifyManager(ctx, dto.Name, time.Now().Add(timeDelta).Unix())
	if err != nil {
		return fmt.Errorf("failed to notify manager due to: %v", err)
	}
	return nil
}

func (s service) NotifyManager(ctx context.Context, gameType string, startTime int64) error {
	u := notifyMangerURL
	dto := NotifyManagerDTO{
		Type:       typeTraining,
		GameType:   gameType,
		Expiration: startTime,
	}
	bytes, err := json.Marshal(&dto)
	if err != nil {
		return fmt.Errorf("failed to marshal data due to: %v", err)
	}
	body := io.NopCloser(strings.NewReader(string(bytes)))
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		log.Printf("failed to do request due to: %v", err)
		return err
	}
	if response == nil {
		//log.Println("response is null")
		return fmt.Errorf("response is null")
	}
	if response.StatusCode != 200 {
		bytes, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		log.Printf(string(bytes))
		return fmt.Errorf("got wrong status code: %d", response.StatusCode)
	}
	return nil
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

func (s service) DeleteCollection(ctx context.Context, dto CollectionDTO) error {
	s.logger.Debug("delete collection")

	err := s.storage.DeleteCollectionByName(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to delete record. error: %w", err)
	}
	return nil
}

// UpdateTime returns new deletion time: + 48 hours
func (s service) UpdateTime(ctx context.Context, tableName string) (int64, error) {
	s.logger.Println("GOT INTO UPDATE TABLE")
	newExpiration := time.Now().Add(timeDelta).Unix()
	dto := CollectionDTO{
		AccessKey: config.GetConfig().Keys.AccessKey,
		Name:      tableName,
	}
	err := s.DeleteCollection(ctx, dto)
	if err != nil {
		return 0, fmt.Errorf("failed to delete collection due to: %v", err)
	}

	err = s.CreateCollection(ctx, dto)
	if err != nil {
		return 0, fmt.Errorf("failed to create collection due to: %v", err)
	}
	return newExpiration, nil
}
