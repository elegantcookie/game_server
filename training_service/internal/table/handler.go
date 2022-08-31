package table

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"training_service/internal/auth"
	"training_service/pkg/logging"
)

var (
	recordsUrl           = "/api/training"
	getAllRecordsUrl     = "/api/training/get/all/"
	getRecordUrl         = "/api/training/get/id/"
	getRecordByUserIDUrl = "/api/training/get/userid"
	getAllCollectionsUrl = "/api/training/collections/get/all"
	collectionsUrl       = "/api/training/collections"
)

type Handler struct {
	Logger          logging.Logger
	TrainingService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, recordsUrl, auth.Middleware(h.CreateRecord))
	router.HandlerFunc(http.MethodPost, getRecordUrl, auth.Middleware(h.GetRecordById))
	router.HandlerFunc(http.MethodPost, getRecordByUserIDUrl, auth.Middleware(h.GetRecordByUserId))
	router.HandlerFunc(http.MethodPost, getAllRecordsUrl, auth.Middleware(h.GetRecords))
	router.HandlerFunc(http.MethodGet, getAllCollectionsUrl, auth.Middleware(h.GetCollectionNames))
	router.HandlerFunc(http.MethodDelete, recordsUrl, auth.Middleware(h.DeleteRecord))
	router.HandlerFunc(http.MethodPatch, recordsUrl, auth.Middleware(h.PartiallyUpdateRecord))
	router.HandlerFunc(http.MethodPost, collectionsUrl, auth.Middleware(h.CreateCollection))
	router.HandlerFunc(http.MethodDelete, collectionsUrl, auth.Middleware(h.DeleteCollectionByName))
}

// Create record
// @Summary Create record endpoint
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body RecordDTO true "create record struct"
// @Tags Records
// @Success 201
// @Failure 400
// @Router /api/training [post]
func (h *Handler) CreateRecord(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE RECORD")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	recordId, err := h.TrainingService.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	tmp := make(map[string]string)
	tmp["record_id"] = recordId
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

// Get record by id
// @Summary Get record by record id
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param data body RecordDTO true "only table_name and id (record_id) required"
// @Tags Records
// @Success 200
// @Failure 400
// @Router /api/training/get/id [post]
func (h *Handler) GetRecordById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET RECORD BY ID")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	user, err := h.TrainingService.GetById(r.Context(), dto)
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

// Get record by user id
// @Summary Get record by user id
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Param userid path string true "User ID"
// @Tags Records
// @Success 200
// @Failure 400
// @Router /api/training/get/{userid} [post]
func (h *Handler) GetRecordByUserId(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET RECORD BY ID")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	user, err := h.TrainingService.GetByUserId(r.Context(), dto)
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

// Get records
// @Summary Get all records of a lobby
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Records
// @Success 200
// @Failure 400
// @Router /api/training/get/all [post]
func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET RECORDS")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	users, err := h.TrainingService.GetAll(r.Context(), dto)
	if err != nil {
		return err
	}

	h.Logger.Println(users)

	userBytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get collection names
// @Summary Get collection names of "training-service" db
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Records
// @Success 200
// @Failure 400
// @Router /api/training/collections/get/all [get]
func (h *Handler) GetCollectionNames(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET RECORDS")
	w.Header().Set("Content-Type", "application/json")

	collectionNames, err := h.TrainingService.GetCollectionNames(context.Background())
	if err != nil {
		return err
	}

	h.Logger.Println(collectionNames)

	userBytes, err := json.Marshal(collectionNames)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Partially update record
// @Summary Partially update record by user id
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Records
// @Success 204
// @Failure 400
// @Router /api/training [patch]
func (h *Handler) PartiallyUpdateRecord(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE USER")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	err := h.TrainingService.Update(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete record
// @Summary Delete record by record id
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Records
// @Success 204
// @Failure 400
// @Router /api/training [delete]
func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE RECORD")
	w.Header().Set("Content-Type", "application/json")

	var dto RecordDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return fmt.Errorf("unable to decode response body due to: %v", err)
	}

	err = h.TrainingService.Delete(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Create collection
// @Summary Create collection endpoint. Needs accept token
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Collections
// @Success 201
// @Failure 400
// @Router /api/training/collections [post]
func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE COLLECTION")
	w.Header().Set("Content-Type", "application/json")

	var dto CollectionDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	ctx := context.WithValue(r.Context(), "Authorization", r.Header.Get("Authorization"))
	err = h.TrainingService.CreateCollection(ctx, dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

// Delete collection
// @Summary Delete collection by collection name(table_name). Needs accept token
// @Accept json
// @Produce json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)

// @Tags Collections
// @Success 204
// @Failure 400
// @Router /api/training/collections [delete]
func (h *Handler) DeleteCollectionByName(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE COLLECTION")
	w.Header().Set("Content-Type", "application/json")

	var dto CollectionDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	err = h.TrainingService.DeleteCollection(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
