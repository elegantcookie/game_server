package ticket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"ticket_service/internal/auth"
	"ticket_service/internal/config"
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
	UseTicket(ctx context.Context, ticketID string) error
	SetFreeTicketStatus(dto FreeTicketStatusDTO) error
	GetFreeTicketStatus() bool
}

func (s service) Create(ctx context.Context, dto TicketDTO) (ticketID string, err error) {
	ticket := Ticket{
		UserID:       dto.UserID,
		IsActive:     true,
		IsGift:       dto.IsGift,
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

	var client http.Client
	addTicketDTO := NewAddTicketDTO(ticketID, &dto)
	bytes, err := json.Marshal(addTicketDTO)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data due to: %v", err)
	}
	reqBody := io.NopCloser(strings.NewReader(string(bytes)))
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, AddTicketByIDURL, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to make request due to: %v", err)
	}
	request.Header.Set("Authorization", dto.JWTToken)
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("failed to do request due to: %v", err)
	}
	if response == nil {
		return "", fmt.Errorf("response is nil")
	}
	if response.StatusCode != 201 {
		if response.StatusCode == 418 {
			bytes, err = io.ReadAll(response.Body)
			if err != nil {
				return "", fmt.Errorf("faileed to read response body due to: %v", err)
			}
			return "", fmt.Errorf("wrong status code: %d. Response text: %s", response.StatusCode, string(bytes))
		}
		log.Println(response.StatusCode)
		return "", fmt.Errorf("wrong status code: %d", response.StatusCode)
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

func (s service) UseTicket(ctx context.Context, ticketID string) error {
	ticket, err := s.GetById(ctx, ticketID)
	if err != nil {
		return err
	}
	if ticket.IsActive == false {
		return fmt.Errorf("ticket is already used")
	}
	ticket.IsActive = false
	err = s.Update(ctx, ticket)
	if err != nil {
		return err
	}
	return nil

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

func (s service) SetFreeTicketStatus(dto FreeTicketStatusDTO) error {
	cfg := config.GetConfig()
	if dto.AccessKey != cfg.Keys.AccessKey {
		return fmt.Errorf("wrong access key")
	}
	cfg.TicketsAvailable = dto.Status
	return nil
}

func (s service) GetFreeTicketStatus() bool {
	return config.GetConfig().TicketsAvailable
}
