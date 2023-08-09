package handlers

import (
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/go-chi/chi"
)

type OrderHandler struct {
	Service order.OrderService
	JwtAuth *middleware.JwtAuthentication
}

func ProvideOrderHandler(service order.OrderService, jwtAuth *middleware.JwtAuthentication) OrderHandler {
	return OrderHandler{Service: service, JwtAuth: jwtAuth}
}

func (h *OrderHandler) Router(r chi.Router) {
	r.Route("/orders", func(r chi.Router) {
		r.Use(h.JwtAuth.Validate)
		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.AdminOnly)
			r.Get("/", h.HandleGetAll)
		})
	})
}

func (h *OrderHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {

}
