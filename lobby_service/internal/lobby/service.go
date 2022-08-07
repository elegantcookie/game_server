package lobby

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"lobby_service/internal/auth"
	"lobby_service/internal/lobby/api"
	"lobby_service/pkg/logging"
	"net/http"
	"strings"
)

var _ Service = &service{}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(storage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: storage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto LobbyDTO) (string, error)
	GetAll(ctx context.Context) ([]Lobby, error)
	GetById(ctx context.Context, id string) (Lobby, error)
	Update(ctx context.Context, dto Lobby) error
	Delete(ctx context.Context, id string) error
	AddUserToLobby(ctx context.Context, dto JoinLobbyDTO) error
}

func (s service) Create(ctx context.Context, dto LobbyDTO) (lobbyID string, err error) {
	s.logger.Debug("check password")
	lobby := Lobby{
		GameType:    dto.GameType,
		MaxPlayers:  dto.MaxPlayers,
		TicketPrice: dto.TicketPrice,
		PrizeSum:    dto.PrizeSum,
		Players:     []Player{},
		StartTime:   dto.StartTime,
		EndTime:     dto.EndTime,
	}
	lobbyID, err = s.storage.Create(ctx, lobby)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return lobbyID, err
		}
		return lobbyID, fmt.Errorf("failed to create lobby. error: %w", err)
	}

	return lobbyID, nil
}

// GetOne Lobby by id
func (s service) GetById(ctx context.Context, id string) (lobby Lobby, err error) {
	lobby, err = s.storage.FindById(ctx, id)

	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return lobby, err
		}
		return lobby, fmt.Errorf("failed to find Lobby by uuid. error: %w", err)
	}
	return lobby, nil
}

func (s service) GetAll(ctx context.Context) ([]Lobby, error) {
	lobbys, err := s.storage.FindAll(ctx)
	if err != nil {
		return lobbys, fmt.Errorf("failed to find lobbys. error: %v", err)
	}
	return lobbys, nil
}

func (s service) Update(ctx context.Context, lobby Lobby) error {
	err := s.storage.Update(ctx, lobby)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to update lobby. error: %w", err)
	}
	return err
}

func getPlayerIndex(id string, players []Player) (int, bool) {
	for i, player := range players {
		if player.ID == id {
			return i, true
		}
	}
	return -1, false
}

func GetUserByID(ctx context.Context, userID string) (dto UpdateUserDTO, err error) {
	_u := fmt.Sprintf("%s%s", GetUsersByIDURL, userID)

	response, err := api.MakeRequest(http.MethodPost, _u, nil)
	if err != nil {
		return dto, err
	}
	defer response.Body.Close()

	if err != nil {
		return dto, fmt.Errorf("failed to do request due to: %v", err)
	}

	err = json.NewDecoder(response.Body).Decode(&dto)
	if err != nil {
		return dto, fmt.Errorf("failed to parse response body due to: %v", err)
	}
	return dto, nil
}

func getGameTickets(gameType string, gameTicketsArr []GameTickets) (index int, found bool) {
	for i := 0; i < len(gameTicketsArr); i++ {
		if gameTicketsArr[i].GameType == gameType {
			return i, true
		}
	}
	return -1, false
}

func UseTicket(ctx context.Context, ticketID string) error {
	url := fmt.Sprintf("%s/%s", UseTicketURL, ticketID)
	response, err := api.MakeRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to do request due to: %v", err)
	}
	if response.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to use ticket due to: %s", string(bytes))
	}
	return nil
}

func UpdateUserTicket(ctx context.Context, gameType string, dto UpdateUserDTO) error {
	var (
		found bool
		i     int
	)
	if i, found = getGameTickets(gameType, dto.Tickets); !found {
		return fmt.Errorf("user has no tickets of game type: %s", gameType)
	}

	if dto.Tickets[i].Amount == 0 {
		return fmt.Errorf("user has no active tickets of game type: %s", gameType)
	}

	dto.Tickets[i].Amount -= 1

	_url := UpdateUserURL
	bytes, err := json.Marshal(&dto)
	if err != nil {
		return err
	}

	_, err = api.MakeRequest(http.MethodPatch, _url, io.NopCloser(strings.NewReader(string(bytes))))
	if err != nil {
		return fmt.Errorf("failed to do request")
	}
	//if response.StatusCode != 200 {
	//	bytes, err := ioutil.ReadAll(response.Body)
	//	if err != nil {
	//		return err
	//	}
	//	return fmt.Errorf("failed to update user due to: %s", string(bytes))
	//}
	return nil
}

func (s service) AddUserToLobby(ctx context.Context, dto JoinLobbyDTO) error {
	s.logger.Println("GOT INTO addUserToLobby")
	lobby, err := s.storage.FindById(ctx, dto.LobbyID)
	if err != nil {
		return err
	}

	if i, found := getPlayerIndex(dto.UserID, lobby.Players); found {
		if lobby.Players[i].Ready == true {
			return fmt.Errorf("user is already ready")
		}
		lobby.Players[i].Ready = true
	} else {
		if lobby.NowPlayers == lobby.MaxPlayers {
			return fmt.Errorf("lobby is full")
		}

		s.logger.Printf("trying to get user by id: %s", dto.UserID)

		user, err := GetUserByID(ctx, dto.UserID)
		if err != nil {
			return err
		}

		s.logger.Println("got user by id: %v", user)

		userDTO := UpdateUserDTO{
			ID:            user.ID,
			Username:      user.Username,
			HasFreeTicket: user.HasFreeTicket,
			Tickets:       user.Tickets,
		}
		err = UpdateUserTicket(ctx, lobby.GameType, userDTO)
		if err != nil {
			return fmt.Errorf("failed to update user ticket due to: %v", err)
		}

		err = UseTicket(ctx, dto.TicketID)
		if err != nil {
			return fmt.Errorf("failed to use ticket due to: %v", err)
		}

		lobby.Players = append(lobby.Players, Player{
			ID:    dto.UserID,
			Ready: false,
		})
		lobby.NowPlayers += 1
	}

	err = s.storage.Update(ctx, lobby)
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
		return fmt.Errorf("failed to delete lobby. error: %w", err)
	}
	return err
}
