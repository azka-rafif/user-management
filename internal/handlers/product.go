package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/product"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/shared/pagination"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type ProductHandler struct {
	Service product.ProductService
	JwtAuth *middleware.JwtAuthentication
}

func ProvideProductHandler(service product.ProductService, jwtAuth *middleware.JwtAuthentication) ProductHandler {
	return ProductHandler{Service: service, JwtAuth: jwtAuth}
}

func (h *ProductHandler) Router(r chi.Router) {
	r.Route("/products", func(r chi.Router) {
		r.Use(h.JwtAuth.Validate)

		r.Group(func(r chi.Router) {
			r.Get("/", h.HandleGetAll)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.AdminOnly)
			r.Post("/", h.HandleCreateProduct)
		})
	})
}

func (h *ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload product.ProductPayload
	err := decoder.Decode(&payload)
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

	res, err := h.Service.Create(payload, userId)

	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, res)
}

func (h *ProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
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
	productTitle := pagination.ParseQueryParams(r, "product_title")
	offset := (page - 1) * limit
	res, err := h.Service.GetAll(limit, offset, sort, field, productTitle)
	totalPage := int(math.Ceil(float64(len(res)) / float64(limit)))

	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithPagination(w, http.StatusOK, res, page, limit, totalPage)
}
