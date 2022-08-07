package lobby

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"lobby_service/internal/auth"
	"lobby_service/pkg/logging"
	"net/http"
)

var (
	lobbiesUrl      = "/api/lobbies"
	getAllLobbysUrl = "/api/lobbies/all"
	lobbyUrl        = "/api/lobbies/id/:id"
	joinLobbyURL    = "/api/lobbies/join"
)

type Handler struct {
	Logger       logging.Logger
	LobbyService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, lobbiesUrl, auth.Middleware(h.CreateLobby))
	router.HandlerFunc(http.MethodPost, lobbyUrl, auth.Middleware(h.GetLobbyById))
	router.HandlerFunc(http.MethodPost, getAllLobbysUrl, auth.Middleware(h.GetLobbys))
	router.HandlerFunc(http.MethodDelete, lobbyUrl, auth.Middleware(h.DeleteLobby))
	router.HandlerFunc(http.MethodPatch, lobbiesUrl, auth.Middleware(h.PartiallyUpdateLobby))
	router.HandlerFunc(http.MethodPost, joinLobbyURL, auth.Middleware(h.JoinLobby))
}

// Create lobby
// @Summary Create lobby endpoint
// @Accept json
// @Produce json
// @Tags Lobbys
// @Success 201
// @Failure 400
// @Router /api/lobbies [post]
func (h *Handler) CreateLobby(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE LOBBY")
	w.Header().Set("Content-Type", "application/json")

	var dto LobbyDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	lobbyID, err := h.LobbyService.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	tmp := make(map[string]string)
	tmp["lobby_id"] = lobbyID
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

// Get lobby by id
// @Summary Get lobby by lobby id
// @Accept json
// @Produce json
// @Tags Lobbys
// @Success 200
// @Failure 400
// @Router /api/lobbies/id/:id [post]
func (h *Handler) GetLobbyById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LOBBY BY ID")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	lobby, err := h.LobbyService.GetById(r.Context(), id)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal lobby")
	lobbyBytes, err := json.Marshal(lobby)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lobbyBytes)
	return nil
}

// Get lobbies
// @Summary Get all lobbies
// @Accept json
// @Produce json
// @Tags Lobbys
// @Success 200
// @Failure 400
// @Router /api/lobbies/all [post]
func (h *Handler) GetLobbys(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LOBBYS")
	w.Header().Set("Content-Type", "application/json")

	lobbies, err := h.LobbyService.GetAll(r.Context())
	if err != nil {
		return err
	}

	h.Logger.Println(lobbies)

	lobbyBytes, err := json.Marshal(lobbies)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lobbyBytes)
	return nil
}

// Partially update lobby
// @Summary Partially update lobby by lobby id
// @Accept json
// @Produce json
// @Tags Lobbys
// @Success 204
// @Failure 400
// @Router /api/lobbies [patch]
func (h *Handler) PartiallyUpdateLobby(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE LOBBY")
	w.Header().Set("Content-Type", "application/json")

	var lobby Lobby
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&lobby); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.LobbyService.Update(r.Context(), lobby)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete lobby
// @Summary Delete lobby by lobby id
// @Accept json
// @Produce json
// @Tags Lobbys
// @Success 204
// @Failure 400
// @Router /api/lobbies/id/:id [delete]
func (h *Handler) DeleteLobby(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE LOBBY")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	err := h.LobbyService.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) JoinLobby(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("JOIN LOBBY")
	w.Header().Set("Content-Type", "application/json")

	var dto JoinLobbyDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	fmt.Printf("%v", dto)
	err := h.LobbyService.AddUserToLobby(context.Background(), dto)
	if err != nil {
		return fmt.Errorf("failed to add user to lobby due to: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
