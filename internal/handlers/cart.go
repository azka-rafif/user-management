package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/shared/pagination"
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
		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.AdminOnly)
			r.Get("/", h.HandleGetAllCarts)
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
// @Router /v1/carts/{cartId} [post]
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
// @Success 200 {object} response.Base{data=cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts/{cartId} [get]
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
// @Success 200 {object} response.Base{data=[]cart.CartItemResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts/{cartId}/items [get]
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
// @Router /v1/carts/{cartId}/checkout [post]
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

// HandleGetAll Gets all carts.
// @Summary Gets all carts.
// @Description This endpoint Gets all carts of users.
// @Tags v1/Cart
// @Security JWTToken
// @Param page query int true "current page number"
// @Param limit query int true "limit of carts per page"
// @Param sort query string false "sort direction"
// @Param field query string false "field to sort by"
// @Produce json
// @Success 200 {object} response.Pagination{data=[]cart.CartResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/carts [get]
func (h *CartHandler) HandleGetAllCarts(w http.ResponseWriter, r *http.Request) {
	page, err := pagination.ConvertToInt(pagination.ParseQueryParams(r, "page"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	limit, err := pagination.ConvertToInt(pagination.ParseQueryParams(r, "limit"))
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	sort := pagination.GetSortDirection(pagination.ParseQueryParams(r, "sort"))
	field := pagination.CheckFieldQuery(pagination.ParseQueryParams(r, "field"), "id")
	offset := (page - 1) * limit
	res, err := h.Service.GetAllCarts(limit, offset, sort, field)
	totalPage := int(math.Ceil(float64(len(res)) / float64(limit)))

	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithPagination(w, http.StatusOK, res, page, limit, totalPage)
}
