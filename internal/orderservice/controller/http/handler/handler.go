package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"orderservice/internal/orderservice/service"
	"time"

	"github.com/go-chi/chi"
)

const reqTimeOut = time.Second * 3

const ErrBadRequest = "ErrBadRequest"
const ErrInternalServerError = "ErrInternalServerError"
const ErrNotFound = "ErrNotFound"

type handler struct {
	timeout time.Duration
	svc     service.OrderService
}

// Constructor
func New(svc service.OrderService) *handler {
	return &handler{
		timeout: reqTimeOut,
		svc:     svc,
	}
}

// Register routes on router
func (h *handler) RegisterRoutes(router *chi.Mux) {

	router.Get("/order", h.Get)
}

// Get order by id

//		@Summary		Show a order
//		@Description	get {object} model.Order by ID
//		@Tags			orders
//		@Accept			json
//		@Produce		json
//	    @Param          id   query      string  true  "Order ID"
//		@Success		200	{object}	model.Order
//		@Failure		400	string	ErrBadRequest
//		@Failure		404	string	ErrNotFound
//		@Failure		500	string	ErrInternalServerError
//		@Router			/order [get]
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	var input InputDTO
	input.Id = r.FormValue("id")
	defer r.Body.Close()

	id, err := input.toModel()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(ErrBadRequest))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()
	orders, err := h.svc.GetById(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(ErrNotFound))
		return
	}

	orderBytes, err := json.Marshal(orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(ErrInternalServerError))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(orderBytes)
}
