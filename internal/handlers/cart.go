package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
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
			r.Get("/items", h.HandleGetCartItems)
			r.Post("/checkout", h.HandleCheckout)
		})
	})
}

// HandleRegister Adds a product into a cart.
// @Summary Adds a product into a users cart.
// @Description This endpoint Creates a cart item and put it into a users cart.
// @Tags v1/Cart
// @Security JWTToken
// @Param cartId path string true "the cart id"
// @Param CartItem body cart.CartItemPayload true "The product to be added to the cart."
// @Produce json
// @Success 201 {object} response.Base{data=cart.CartItemResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/{cartId} [post]
func (h *CartHandler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "cartId")
	cartId, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload cart.CartItemPayload
	err = decoder.Decode(&payload)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = shared.GetValidator().Struct(payload)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userId, err := uuid.FromString(claims.UserId)

	if err != nil {
		response.WithError(w, err)
		return
	}

	res, err := h.Service.AddToCart(payload, userId, cartId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, res)
}

// HandleGetCart Gets a Cart.
// @Summary Gets a users cart.
// @Description This endpoint Gets a users cart.
// @Tags v1/Cart
// @Security JWTToken
// @Param cartId path string true "the cart id"
// @Produce json
// @Success 201 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/{cartId} [get]
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

// HandleGetCartItems Gets only the carts items.
// @Summary Gets the cart items.
// @Description This endpoint only gets the carts items.
// @Tags v1/Cart
// @Security JWTToken
// @Param cartId path string true "the cart id"
// @Produce json
// @Success 201 {object} response.Base{data=[]cart.CartItemResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/{cartId}/items [get]
func (h *CartHandler) HandleGetCartItems(w http.ResponseWriter, r *http.Request) {
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
	response.WithJSON(w, http.StatusOK, res.CartItems)
}

// HandleCheckout checkout a list of cart items.
// @Summary checkout a list of cart items.
// @Description This endpoint checkout the list of cart item ids given.
// @Tags v1/Cart
// @Security JWTToken
// @Param cartId path string true "the cart id"
// @Param CartItemIds body cart.CheckoutPayload true "The items to be checked out"
// @Produce json
// @Success 201 {object} response.Base{data=order.OrderResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/cart/{cartId}/checkout [post]
func (h *CartHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "cartId")
	cartId, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload cart.CheckoutPayload
	err = decoder.Decode(&payload)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	err = shared.GetValidator().Struct(payload)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userId, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	res, err := h.Service.Checkout(payload, cartId, userId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, res)
}
