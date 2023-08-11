package handlers

import (
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/order"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/shared/pagination"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
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
		r.Get("/", h.HandleGetAll)
		r.Route("/{orderId}", func(r chi.Router) {
			r.Get("/", h.HandleGetByID)
			r.Delete("/", h.HandleCancel)
		})
	})
}

// HandleGetAll Gets all orders.
// @Summary Gets all orders.
// @Description This endpoint Gets all orders of users if the current user is admin if not it gets only the active users orders.
// @Tags v1/Order
// @Security JWTToken
// @Param page query int true "current page number"
// @Param limit query int true "limit of orders per page"
// @Param sort query string false "sort direction"
// @Param field query string false "field to sort by"
// @Param status query string false "filter by order status"
// @Produce json
// @Success 200 {object} response.Base{data=[]order.OrderResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/orders [get]
func (h *OrderHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	pg, err := pagination.GetPagination(r)
	if err != nil {
		response.WithError(w, err)
		return
	}
	status := pagination.ParseQueryParams(r, "status")
	cancelled := pagination.GetCancelled(pagination.ParseQueryParams(r, "cancelled"))
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
	res, err := h.Service.GetAll(pg.Limit, pg.Offset, pg.Sort, pg.Field, status, userId, claims.Role, cancelled)
	totalPage := pg.GetTotalPages(res)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithPagination(w, http.StatusOK, res, pg.Page, pg.Limit, totalPage)
}

// HandleCancel cancel an order.
// @Summary Cancels an Order.
// @Description This endpoint cancels an active order.
// @Tags v1/Order
// @Security JWTToken
// @Param orderId path string true "the order id"
// @Produce json
// @Success 200 {object} response.Base{data=order.OrderResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/orders/{orderId} [delete]
func (h *OrderHandler) HandleCancel(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "orderId")
	id, err := uuid.FromString(idString)
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
	res, err := h.Service.CancelOrder(id, userId, claims.Role)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}

// Handleget get an order by id.
// @Summary gets an Order by id.
// @Description This endpoint gets an order by id.
// @Tags v1/Order
// @Security JWTToken
// @Param orderId path string true "the order id"
// @Produce json
// @Success 200 {object} response.Base{data=order.OrderResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/orders/{orderId} [delete]
func (h *OrderHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "orderId")
	id, err := uuid.FromString(idString)
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
	res, err := h.Service.GetByID(id, userId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}
