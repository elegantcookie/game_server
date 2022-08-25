package snake

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"snake_service/internal/auth"
	"snake_service/pkg/logging"
)

var (
	gameServersUrl  = "/api/snake/"
	getAllSnakesUrl = "/api/snake/all/"
	gameServerIDUrl = "/api/snake/id/:id"
	sendResultURL   = "/api/snake/res/"
	getStatusURL    = "/api/snake/status/:id"
)

type Handler struct {
	Logger      logging.Logger
	GameService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, gameServersUrl, auth.Middleware(h.CreateGS))
	router.HandlerFunc(http.MethodPost, gameServerIDUrl, auth.Middleware(h.GetGSById))
	router.HandlerFunc(http.MethodPost, getAllSnakesUrl, auth.Middleware(h.GetGameServers))
	router.HandlerFunc(http.MethodDelete, gameServerIDUrl, auth.Middleware(h.DeleteGS))
	router.HandlerFunc(http.MethodPut, gameServersUrl, auth.Middleware(h.PartiallyUpdateGS))
	router.HandlerFunc(http.MethodPost, sendResultURL, auth.Middleware(h.SendResult))
	router.HandlerFunc(http.MethodPost, getStatusURL, auth.Middleware(h.GetGameStatus))
}

// Create game server
// @Summary Create game server endpoint
// @Accept json
// @Produce json
// @Tags Snakes
// @Success 201
// @Failure 400
// @Router /api/snakes [post]
func (h *Handler) CreateGS(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE GAME SERVER")
	w.Header().Set("Content-Type", "application/json")

	var dto SnakeDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	log.Println(dto)
	snakeID, err := h.GameService.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	tmp := make(map[string]string)
	tmp["snake_id"] = snakeID
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

// Get game server by id
// @Summary Get game server by game server id
// @Accept json
// @Produce json
// @Tags Snakes
// @Success 200
// @Failure 400
// @Router /api/snakes/get/id [post]
func (h *Handler) GetGSById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET GAME SERVER BY ID")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	gsID := params.ByName("id")
	user, err := h.GameService.GetById(r.Context(), gsID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal game server")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall game server. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get snakes
// @Summary Get all snakes
// @Accept json
// @Produce json
// @Tags Snakes
// @Success 200
// @Failure 400
// @Router /api/snakes/get/all [post]
func (h *Handler) GetGameServers(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET GAME SERVERS")
	w.Header().Set("Content-Type", "application/json")

	snakes, err := h.GameService.GetAll(r.Context())
	if err != nil {
		return err
	}

	h.Logger.Println(snakes)

	userBytes, err := json.Marshal(snakes)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Partially update game server
// @Summary Partially update game server by user id
// @Accept json
// @Produce json
// @Tags Snakes
// @Success 204
// @Failure 400
// @Router /api/snakes [patch]
func (h *Handler) PartiallyUpdateGS(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE GAME SERVER")
	w.Header().Set("Content-Type", "application/json")

	var snake Snake
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&snake); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.GameService.Update(r.Context(), snake)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete game server
// @Summary Delete game server by game server id
// @Accept json
// @Produce json
// @Tags Snakes
// @Success 204
// @Failure 400
// @Router /api/snakes [delete]
func (h *Handler) DeleteGS(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE GAME SERVER")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	gsID := params.ByName("id")

	err := h.GameService.Delete(r.Context(), gsID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) SendResult(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("SEND RESULT")
	w.Header().Set("Content-Type", "application/json")

	var dto SendResultDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return fmt.Errorf("failed to decode response body due to: %v", err)
	}
	err = h.GameService.SendResult(r.Context(), dto)
	if err != nil {
		return fmt.Errorf("failed to send result due to: %v", err)
	}
	return nil

}

func (h *Handler) GetGameStatus(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET STATUS")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")
	status, err := h.GameService.GetGameStatus(r.Context(), id)
	if err != nil {
		return fmt.Errorf("failed to get game status due to: %v", err)
	}
	tmp := map[string]int{"status": status}
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %v", err)
	}
	w.Write(bytes)
	return nil
}
