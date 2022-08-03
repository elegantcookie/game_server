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
	ticketsUrl         = "/api/tickets"
	getAllTicketsUrl   = "/api/tickets/get/all"
	getTicketUrl       = "/api/tickets/get/id"
	getTicketStatusURL = "/api/tickets/get/status/id"
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
}

// Create ticket
// @Summary Create ticket endpoint
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

// Get ticket by id
// @Summary Get ticket by ticket id
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

	h.Logger.Debug("marshal ticket")
	userBytes, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshall ticket. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
	return nil
}

// Get ticket status
// @Summary Get ticket status by ticket id
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
		return fmt.Errorf("failed to marshall ticket. error: %w", err)
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

// Partially update ticket
// @Summary Partially update ticket by user id
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
	h.Logger.Printf("Ticket: %v", ticket)
	err := h.TicketService.Update(r.Context(), ticket)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// Delete ticket
// @Summary Delete ticket by ticket id
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
