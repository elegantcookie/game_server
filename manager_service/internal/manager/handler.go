package manager

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"manager_service/internal/apperror"
	"manager_service/pkg/logging"
	"net/http"
)

type Handler struct {
	Logger         logging.Logger
	ManagerService Service
}

var (
	mainURL      = "/api/manager/"
	getRLSURL    = "/api/manager/all"
	getRLURL     = "/api/manager/id/:id"
	deleteAllURL = "/api/manager/del/all"
)

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, mainURL, apperror.NoAuthMiddleware(h.Create))
	router.HandlerFunc(http.MethodPost, getRLSURL, apperror.Middleware(h.GetLRS))
	router.HandlerFunc(http.MethodPost, getRLURL, apperror.Middleware(h.GetLRById))
	router.HandlerFunc(http.MethodDelete, getRLURL, apperror.Middleware(h.DeleteLR))
	router.HandlerFunc(http.MethodDelete, deleteAllURL, apperror.Middleware(h.DeleteAll))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) error {
	log.Println("CREATE LR RECORD")
	w.Header().Set("Content-Type", "application/json")
	var dto LobbyRecordDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return fmt.Errorf("failed to decode due to: %v", err)
	}
	fmt.Println(dto)
	id, err := h.ManagerService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	tmp := map[string]string{"id": id}
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil

}

func (h *Handler) GetLRS(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LRS")
	w.Header().Set("Content-Type", "application/json")

	lobbies, err := h.ManagerService.GetAll(r.Context())
	if err != nil {
		return err
	}

	lobbyBytes, err := json.Marshal(lobbies)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(lobbyBytes)
	return nil
}

func (h *Handler) GetLRById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET LR BY ID")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	lobby, err := h.ManagerService.GetById(r.Context(), id)
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

func (h *Handler) DeleteLR(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE LR")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	err := h.ManagerService.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE ALL")
	w.Header().Set("Content-Type", "application/json")

	err := h.ManagerService.DeleteAll(r.Context())
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
