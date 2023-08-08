package handlers

import (
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type CartHandler struct {
	Service cart.CartService
	JwtAuth *middleware.JwtAuthentication
}

func ProvideCartHandler(service cart.CartService, jwtAuth *middleware.JwtAuthentication) CartHandler {
	return CartHandler{Service: service, JwtAuth: jwtAuth}
}

func (h *CartHandler) Router(r chi.Router) {
	r.Route("/carts", func(r chi.Router) {
		r.Use(h.JwtAuth.Validate)
		r.Route("/{cartId}", func(r chi.Router) {
			r.Use(h.JwtAuth.CartAccess)
			r.Get("/", h.HandleGetCart)
			r.Post("/", h.HandleAddToCart)
		})
	})
}

func (h *CartHandler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {

}

func (h *CartHandler) HandleGetCart(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "cartId")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	res, err := h.Service.GetCart(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}
