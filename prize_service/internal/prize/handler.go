package prize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"prize_service/internal/auth"
	"prize_service/pkg/logging"
)

var (
	prizesUrl       = "/api/prizes"
	getAllPrizesUrl = "/api/prizes/all"
	prizeUrl        = "/api/prizes/id/:id"
)

type Handler struct {
	Logger       logging.Logger
	PrizeService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, prizesUrl, auth.Middleware(h.CreatePrize))
	router.HandlerFunc(http.MethodPost, prizeUrl, auth.Middleware(h.GetPrizeById))
	router.HandlerFunc(http.MethodPost, getAllPrizesUrl, auth.Middleware(h.GetPrizes))
	router.HandlerFunc(http.MethodDelete, prizeUrl, auth.Middleware(h.DeletePrize))
	router.HandlerFunc(http.MethodPatch, prizesUrl, auth.Middleware(h.PartiallyUpdatePrize))
}

// Create prize
// @Summary Create prize endpoint
// @Accept json
// @Produce json
// @Tags Prizes
// @Success 201
// @Failure 400
// @Router /api/prizes [post]
func (h *Handler) CreatePrize(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE PRIZE")
	w.Header().Set("Content-Type", "application/json")

	var dto PrizeDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	prizeID, err := h.PrizeService.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	tmp := make(map[string]string)
	tmp["prize_id"] = prizeID
	bytes, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

// Get prize by id
// @Summary Get prize by prize id
// @Accept json
// @Produce json
// @Tags Prizes
// @Success 200
// @Failure 400
// @Router /api/prizes/id/:id [post]
func (h *Handler) GetPrizeById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET PRIZE BY ID")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	prize, err := h.PrizeService.GetById(r.Context(), id)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal prize")
	prizeBytes, err := json.Marshal(prize)
	if err != nil {
		return fmt.Errorf("failed to marshall prize. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(prizeBytes)
	return nil
}

// Get prizes
// @Summary Get all prizes
// @Accept json
// @Produce json
// @Tags Prizes
// @Success 200
// @Failure 400
// @Router /api/prizes/all [post]
func (h *Handler) GetPrizes(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET PRIZES")
	w.Header().Set("Content-Type", "application/json")

	prizes, err := h.PrizeService.GetAll(r.Context())
	if err != nil {
		return err
	}

	h.Logger.Println(prizes)

	prizeBytes, err := json.Marshal(prizes)
	if err != nil {
		return fmt.Errorf("failed to marshall prize. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(prizeBytes)
	return nil
}

// Partially update prize
// @Summary Partially update prize by prize id
// @Accept json
// @Produce json
// @Tags Prizes
// @Success 204
// @Failure 400
// @Router /api/prizes [patch]
func (h *Handler) PartiallyUpdatePrize(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE PRIZE")
	w.Header().Set("Content-Type", "application/json")

	var prize Prize
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&prize); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.PrizeService.Update(r.Context(), prize)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete prize
// @Summary Delete prize by prize id
// @Accept json
// @Produce json
// @Tags Prizes
// @Success 204
// @Failure 400
// @Router /api/prizes/id/:id [delete]
func (h *Handler) DeletePrize(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE PRIZE")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	id := params.ByName("id")

	err := h.PrizeService.Delete(r.Context(), id)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}
