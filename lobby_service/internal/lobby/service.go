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
	"log"
	"net/http"
	"strings"
	"time"
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
	DeleteAll(ctx context.Context) error
	AddUserToLobby(ctx context.Context, dto JoinLobbyDTO) error
	GetLobbyIDByParams(ctx context.Context, params Params) (string, error)
	UpdateLobbyTime(ctx context.Context, utdto UpdateTimeDTO) (int64, error)
}

func NotifyManager(ctx context.Context, jwtToken, lobbyID string, startTime int64) error {
	log.Printf("JWT TOKEN: %v", jwtToken)
	u := "http://localhost:10007/api/manager/"

	body := io.NopCloser(strings.NewReader(fmt.Sprintf(`
{
	"lobby_id": "%s",
	"expiration": %d
}`, lobbyID, startTime)))
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	log.Println(jwtToken)
	request.Header.Add("Authorization", jwtToken)
	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		log.Printf("failed to do request due to: %v", err)
		return err
	}
	if response == nil {
		log.Println("response is null")
		return fmt.Errorf("response is null")
	}
	if response.StatusCode != 200 {
		log.Printf("got wrong status code: %d", response.StatusCode)

		return fmt.Errorf("got wrong status code: %d", response.StatusCode)
	}
	return nil
}

func (s service) Create(ctx context.Context, dto LobbyDTO) (lobbyID string, err error) {
	s.logger.Debug("CREATE LOBBY SERVICE")
	lobby := Lobby{
		GameType:    dto.GameType,
		MaxPlayers:  dto.MaxPlayers,
		TicketPrice: dto.TicketPrice,
		PrizeSum:    dto.PrizeSum,
		Players:     []Player{},
		StartTime:   dto.StartTime,
		EndTime:     dto.EndTime,
	}

	var start int
	now := time.Now()
	hour := now.Hour()
	if hour%2 == 0 {
		start = hour + 2
	} else {
		start = hour + 1
	}

	// If start time and is not in a payload then it sets start time as next even hour
	if dto.StartTime == 0 {
		startTime := time.Date(now.Year(), now.Month(), now.Day(), start, 0, 0, 0, now.Location()).Unix()
		lobby.StartTime = startTime
	}

	lobby.EndTime = lobby.StartTime + TwoHours

	log.Printf("CREATING LOBBY WITH START TIME: %v", lobby.StartTime)

	lobbyID, err = s.storage.Create(ctx, lobby)
	if err != nil {
		if errors.Is(err, auth.ErrNotFound) {
			return lobbyID, err
		}
		return lobbyID, fmt.Errorf("failed to create lobby. error: %w", err)
	}

	err = NotifyManager(ctx, dto.JWTToken, lobbyID, lobby.StartTime)
	if err != nil {
		return "", fmt.Errorf("failed to notify manager due to: %v", err)
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

func (s service) DeleteAll(ctx context.Context) error {
	err := s.storage.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete lr due to: %v", err)
	}
	return nil
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

	response, err := api.MakeRequestWithContext(ctx, http.MethodPost, _u, nil)
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
	response, err := api.MakeRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to do request due to: %v", err)
	}
	if response.StatusCode != 200 {
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to use snake due to: %s", string(bytes))
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

	response, err := api.MakeRequestWithContext(ctx, http.MethodPost, _url, io.NopCloser(strings.NewReader(string(bytes))))
	if err != nil {
		return fmt.Errorf("failed to do request due to: %v", err)
	}
	if response.StatusCode != 204 {
		return fmt.Errorf("user service returned wrong status code: %d", response.StatusCode)
	}
	return nil
}

// RecreateLobby creates lobby using dto and deletes lobby by lobbyID
func (s service) RecreateLobby(ctx context.Context, dto LobbyDTO, lobbyID string) error {
	_, err := s.Create(ctx, dto)
	if err != nil {
		return err
	}
	err = s.Delete(ctx, lobbyID)
	if err != nil {
		return err
	}
	return nil
}

func (s service) CreateSnakeGS(lobby Lobby) error {
	log.Println("CREATE GAME SERVER")
	log.Printf("%v", lobby)
	ids := GetPlayersIDS(lobby)
	var dto CreateGSDTO
	dto.Players = ids
	dto.StartTime = lobby.StartTime
	dto.EndTime = lobby.EndTime

	bytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal data due to: %v", err)
	}
	request, err := http.NewRequest(http.MethodPost, createGameServerURL, io.NopCloser(strings.NewReader(string(bytes))))
	if err != nil {
		return fmt.Errorf("failed to create new request due to: %v", err)
	}

	request.Header.Add("Authorization", lobby.JWTToken)

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to do request due to: %v", err)
	}
	if response == nil {
		return fmt.Errorf("response is nil")
	}
	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed status code: %d", response.StatusCode)
	}
	return nil
}

// AddUserToLobby needs check if user has such snake.
// Or check must be on client
func (s service) AddUserToLobby(ctx context.Context, dto JoinLobbyDTO) error {
	s.logger.Println("GOT INTO addUserToLobby")
	lobby, err := s.storage.FindById(ctx, dto.LobbyID)
	if err != nil {
		return err
	}

	var (
		i     int
		found bool
	)
	// If user is in players list then set his status ready
	if i, found = getPlayerIndex(dto.UserID, lobby.Players); found {
		if lobby.Players[i].Ready == true {
			return fmt.Errorf("user is already ready")
		}
		lobby.Players[i].Ready = true
	} else {
		// If user not in players list
		// If lobby is full raises error
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
			return fmt.Errorf("failed to update user snake due to: %v", err)
		}

		err = UseTicket(ctx, dto.TicketID)
		if err != nil {
			return fmt.Errorf("failed to use snake due to: %v", err)
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

	// If user is added successfully and lobby got full adds new lobby of the same type
	if lobby.NowPlayers == lobby.MaxPlayers && !found {
		lobbyDTO := LobbyDTO{
			GameType:    lobby.GameType,
			MaxPlayers:  lobby.MaxPlayers,
			NowPlayers:  0,
			TicketPrice: lobby.TicketPrice,
			PrizeSum:    lobby.PrizeSum,
			PrizeType:   lobby.PrizeType,
			StartTime:   lobby.StartTime,
			EndTime:     lobby.EndTime,
			JWTToken:    dto.JWTToken,
		}

		// create new lobby with same params
		_, err = s.Create(ctx, lobbyDTO)
		if err != nil {
			return err
		}
		lobby.JWTToken = dto.JWTToken

		switch lobby.GameType {
		case "snake":
			{
				err := s.CreateSnakeGS(lobby)
				if err != nil {
					return fmt.Errorf("failed to create snake game server due to: %v", err)
				}
			}
		}

		// ????
		// wait for start time to make sure all players are ready
		// start game service after 2 minutes from the start

		// delete lobby
		err = s.Delete(ctx, lobby.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s service) GetLobbyIDByParams(ctx context.Context, params Params) (lobbyID string, err error) {
	lobbyID, err = s.storage.FindByParams(ctx, params.GameType, params.MaxPlayers, params.PrizeSum)
	if err != nil {
		return "", err
	}
	return lobbyID, nil
}

// UpdateLobbyTime DONE
func (s service) UpdateLobbyTime(ctx context.Context, utdto UpdateTimeDTO) (int64, error) {
	lobby, err := s.storage.FindById(ctx, utdto.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to find lobby by id due to: %v", err)
	}

	dto := LobbyDTO{
		GameType:    lobby.GameType,
		MaxPlayers:  lobby.MaxPlayers,
		NowPlayers:  0,
		TicketPrice: lobby.TicketPrice,
		PrizeSum:    lobby.PrizeSum,
		PrizeType:   lobby.PrizeType,
		StartTime:   lobby.StartTime + 24*OneHour,
		EndTime:     lobby.EndTime + 24*OneHour,
		JWTToken:    utdto.JWTToken,
	}
	_, err = s.Create(ctx, dto)
	if err != nil {
		return 0, fmt.Errorf("failed to create new lobby with same params due to: %v", err)
	}

	hour := time.Unix(lobby.StartTime, 0).Hour()
	if hour == 2 {
		lobby.StartTime += 14 * OneHour
	} else {
		lobby.StartTime += OneHour
	}
	lobby.EndTime = lobby.StartTime + 2*OneHour

	err = s.storage.Update(ctx, lobby)
	if err != nil {
		return 0, fmt.Errorf("failed to update lobby time due to: %v", err)
	}
	return lobby.StartTime, nil
}
