package ticket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"ticket_service/internal/auth"
	"ticket_service/pkg/logging"
)

var (
	ticketsUrl             = "/api/tickets"
	getAllTicketsUrl       = "/api/tickets/get/all"
	getTicketUrl           = "/api/tickets/get/id"
	getTicketStatusURL     = "/api/tickets/get/status/id"
	useTicketURL           = "/api/tickets/use/:id"
	getFreeTicketStatusURL = "/api/tickets/free/get/status"
	setFreeTicketStatusURL = "/api/tickets/free/set/status"
)

type Handler struct {
	Logger        logging.Logger
	TicketService Service
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, ticketsUrl, auth.Middleware(h.CreateTicket))
	router.HandlerFunc(http.MethodPost, getTicketUrl, auth.Middleware(h.GetTicketById))
	router.HandlerFunc(http.MethodPost, getAllTicketsUrl, auth.Middleware(h.GetTickets))
	router.HandlerFunc(http.MethodDelete, ticketsUrl, auth.Middleware(h.DeleteTicket))
	router.HandlerFunc(http.MethodPatch, ticketsUrl, auth.Middleware(h.PartiallyUpdateTicket))
	router.HandlerFunc(http.MethodPost, getTicketStatusURL, auth.Middleware(h.GetTicketStatusById))
	router.HandlerFunc(http.MethodPost, setFreeTicketStatusURL, auth.Middleware(h.SetFreeTicketStatus))
	router.HandlerFunc(http.MethodPost, getFreeTicketStatusURL, auth.Middleware(h.GetFreeTicketStatus))
	router.HandlerFunc(http.MethodPost, useTicketURL, auth.NoAuthMiddleware(h.UseTicket))
}

// Create lobby
// @Summary Create lobby endpoint
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 201
// @Failure 400
// @Router /api/tickets [post]
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("POST CREATE TICKET")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	ticketID, err := h.TicketService.Create(context.Background(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	tmp := make(map[string]string)
	tmp["ticket_id"] = ticketID
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
// @Tags Tickets
// @Success 200
// @Failure 400
// @Router /api/tickets/get/id [post]
func (h *Handler) GetTicketById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET TICKET BY ID")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketIDDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	user, err := h.TicketService.GetById(r.Context(), dto.ID)
	if err != nil {
		return err
	}

	h.Logger.Debug("marshal lobby")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get lobby status
// @Summary Get lobby status by lobby id
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 200
// @Failure 400
// @Router /api/tickets/get/status/id [post]
func (h *Handler) GetTicketStatusById(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET TICKET STATUS BY ID")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketIDDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}

	user, err := h.TicketService.GetById(r.Context(), dto.ID)
	if err != nil {
		return err
	}

	ticketStatus := make(map[string]bool)
	ticketStatus["is_active"] = user.IsActive

	userBytes, err := json.Marshal(ticketStatus)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get tickets
// @Summary Get all tickets
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 200
// @Failure 400
// @Router /api/tickets/get/all [post]
func (h *Handler) GetTickets(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GET TICKETS")
	w.Header().Set("Content-Type", "application/json")

	tickets, err := h.TicketService.GetAll(r.Context())
	if err != nil {
		return err
	}

	h.Logger.Println(tickets)

	userBytes, err := json.Marshal(tickets)
	if err != nil {
		return fmt.Errorf("failed to marshall user. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Partially update lobby
// @Summary Partially update lobby by user id
// @Accept json
// @Produce json
// @Tags Tickets
// @Success 204
// @Failure 400
// @Router /api/tickets [patch]
func (h *Handler) PartiallyUpdateTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("PARTIALLY UPDATE TICKET")
	w.Header().Set("Content-Type", "application/json")

	var ticket Ticket
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.TicketService.Update(r.Context(), ticket)
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
// @Tags Tickets
// @Success 204
// @Failure 400
// @Router /api/tickets [delete]
func (h *Handler) DeleteTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("DELETE TICKET")
	w.Header().Set("Content-Type", "application/json")

	var dto TicketIDDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return fmt.Errorf("unable to decode response body due to: %v", err)
	}

	err = h.TicketService.Delete(r.Context(), dto.ID)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *Handler) UseTicket(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("USE TICKET")
	w.Header().Set("Content-Type", "application/json")

	params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
	userID := params.ByName("id")
	err := h.TicketService.UseTicket(r.Context(), userID)
	if err != nil {
		return fmt.Errorf("failed to use quiz due to: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// Set free lobby status
// @Summary Set free lobby status endpoint. Requires authorization and access key
// @Accept json
// @Produce json
// @Tags FreeTickets
// @Success 200
// @Failure 400
// @Router /api/tickets/free/set/status [post]
func (h *Handler) SetFreeTicketStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var dto FreeTicketStatusDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		return auth.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err = h.TicketService.SetFreeTicketStatus(dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	return nil
}

// Get free lobby status
// @Summary Get free lobby status endpoint. Requires authorization
// @Accept json
// @Produce json
// @Tags FreeTickets
// @Success 200
// @Failure 400
// @Router /api/tickets/free/get/status [post]
func (h *Handler) GetFreeTicketStatus(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	status := h.TicketService.GetFreeTicketStatus()
	statusMap := make(map[string]bool)
	statusMap["tickets_available"] = status

	userBytes, err := json.Marshal(statusMap)
	if err != nil {
		return fmt.Errorf("failed to marshall lobby status. error: %w", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}
