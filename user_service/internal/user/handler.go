package user

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"user_service/internal/apperror"
	"user_service/pkg/logging"
)

const (
	usersURL    = "/api/users"
	authUrl     = "/api/users/auth"
	userIdURL   = "/api/users/id/:id"
	usernameURL = "/api/users/username/:username"
	ticketsURL  = "/api/users/tickets"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, usersURL, apperror.AuthMiddleware(h.GetUsers))
	router.HandlerFunc(http.MethodPost, userIdURL, apperror.AuthMiddleware(h.GetUserById))
	router.HandlerFunc(http.MethodPost, usernameURL, apperror.AuthMiddleware(h.GetUserByUsername))
	router.HandlerFunc(http.MethodPost, authUrl, apperror.Middleware(h.GetUserByUsernameAndPassword))
	router.HandlerFunc(http.MethodPatch, usersURL, apperror.AuthMiddleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userIdURL, apperror.AuthMiddleware(h.DeleteUser))
	router.HandlerFunc(http.MethodPut, ticketsURL, apperror.AuthMiddleware(h.AddTicket))
	router.HandlerFunc(http.MethodDelete, ticketsURL, apperror.AuthMiddleware(h.DeleteTicket))
}

// Get user by id
// @Summary Get user by user id
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/users/id/:id [post]
func (h *Handler) GetUserById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER BY ID")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("id")

	user, err := h.UserService.GetById(r.Context(), userUUID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get user by username
// @Summary Get user by username endpoint
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/users/username/:username [post]
func (h *Handler) GetUserByUsername(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER BY USERNAME")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get username from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	username := params.ByName("username")

	user, err := h.UserService.GetByUsername(r.Context(), username)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get user by username and password
// @Summary Get user by username and password. Needs for authorization
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/users/auth [post]
func (h *Handler) GetUserByUsernameAndPassword(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get username from context")

	var dto CreateUserDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return fmt.Errorf("unable to decode response body due to: %v", err)
	}

	user, err := h.UserService.GetByUsernameAndPassword(r.Context(), dto.Username, dto.Password)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get all users
// @Summary Get all users endpoint
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USERS")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")

	users, err := h.UserService.GetAll(r.Context())
	if err != nil {
		return err
	}

	h.Logger.Println(users)

	h.Logger.Debug("marshal user")
	userBytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Create user
// @Summary Create user endpoint
// @Accept json
// @Produce json
// @Tags Users
// @Success 201
// @Failure 400
// @Router /api/users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CREATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode create user dto")
	var crUser CreateUserDTO
	defer r.Body.Close()
	//log.Printf("r.body: %v", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&crUser); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}

	userUUID, err := h.UserService.Create(r.Context(), crUser)
	if err != nil {
		return err
	}
	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersURL, userUUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

// Partially update user by id
// @Summary Partially update user endpoint
// @Accept json
// @Produce json
// @Tags Users
// @Success 204
// @Failure 400
// @Router /api/users [patch]
func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("decode update user dto")
	var dto UpdateUserDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.UserService.Update(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete user by id
// @Summary Delete user by id endpoint
// @Accept json
// @Produce json
// @Tags Users
// @Success 204
// @Failure 400
// @Router /api/users [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("id")

	err := h.UserService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Add ticket
// @Summary Add ticket by user id and ticket id
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 201
// @Failure 400
// @Router /api/users/tickets [put]
func (h *Handler) AddTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST ADD TICKET")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err = h.UserService.AddTicket(r.Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to add ticket due to: %v", err)
	}
	w.WriteHeader(http.StatusCreated)

	return nil
}

// Delete ticket
// @Summary Delete ticket by user id and ticket id
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 204
// @Failure 400
// @Router /api/users/tickets [delete]
func (h *Handler) DeleteTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE TICKET")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err = h.UserService.DeleteTicket(r.Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to delete ticket due to: %v", err)
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
