package handlers

import (
	"math"
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
			r.Delete("/", h.HandleCancel)
		})
	})
}

func (h *OrderHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
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
	status := pagination.ParseQueryParams(r, "status")
	offset := (page - 1) * limit
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
	res, err := h.Service.GetAll(limit, offset, sort, field, status, userId, claims.Role)
	totalPage := int(math.Ceil(float64(len(res)) / float64(limit)))
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithPagination(w, http.StatusOK, res, page, limit, totalPage)

}

func (h *OrderHandler) HandleCancel(w http.ResponseWriter, r *http.Request) {

}
