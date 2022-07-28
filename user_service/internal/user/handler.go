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
	usersURL = "/api/users"
	userURL  = "/api/users/:id"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetUsers))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUser))
	//router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("id")

	user, err := h.UserService.GetOne(r.Context(), userUUID)
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

//func (h *Handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
//	h.Logger.Info("PARTIALLY UPDATE USER")
//	w.Header().Set("Content-Type", "application/json")
//
//	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
//	userID := params.ByName("id")
//
//	h.Logger.Debug("decode update user dto")
//	var updUser UpdateUserDTO
//	defer r.Body.Close()
//	if err := json.NewDecoder(r.Body).Decode(&updUser); err != nil {
//		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
//	}
//	updUser.ID = userID
//
//	err := h.UserService.Update(r.Context(), updUser)
//	if err != nil {
//		return err
//	}
//	w.WriteHeader(http.StatusNoContent)
//
//	return nil
//}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE USER")
	w.Header().Set("Content-Type", "application/json")

	h.Logger.Debug("get uuid from context")
	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userUUID := params.ByName("uuid")

	err := h.UserService.Delete(r.Context(), userUUID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
