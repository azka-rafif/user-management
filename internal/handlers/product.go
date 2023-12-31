package handlers

import (
	"encoding/json"
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

// HandleCheckout Insert a product.
// @Summary create a new product.
// @Description This endpoint inserts a new product.
// @Tags v1/Product
// @Security JWTToken
// @Param CartItemIds body product.ProductPayload true "product to be created"
// @Produce json
// @Success 201 {object} response.Base{data=product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/products [post]
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

// HandleGetAll Gets all products.
// @Summary Gets all products.
// @Description This endpoint Gets all products available.
// @Tags v1/Product
// @Security JWTToken
// @Param page query int true "current page number"
// @Param limit query int true "limit of products per page"
// @Param sort query string false "sort direction"
// @Param field query string false "field to sort by"
// @Param product_title query string false "filter by product name"
// @Produce json
// @Success 200 {object} response.Base{data=[]product.ProductResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/products [get]
func (h *ProductHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	pg, err := pagination.GetPagination(r)
	if err != nil {
		response.WithError(w, err)
		return
	}
	productTitle := pagination.ParseQueryParams(r, "product_title")
	res, err := h.Service.GetAll(pg.Limit, pg.Offset, pg.Sort, pg.Field, productTitle)
	totalPage := pg.GetTotalPages(res)

	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithPagination(w, http.StatusOK, res, pg.Page, pg.Limit, totalPage)
}
