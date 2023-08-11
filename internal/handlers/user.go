package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/user"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type UserHandler struct {
	Service user.UserService
	jwtAuth *middleware.JwtAuthentication
}

func ProvideUserHandler(service user.UserService) UserHandler {
	return UserHandler{Service: service}
}

func (h *UserHandler) Router(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(h.jwtAuth.Validate)

		r.Group(func(r chi.Router) {
			r.Get("/", h.HandleGetUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(h.jwtAuth.IsUser)
			r.Route("/{userId}", func(r chi.Router) {
				r.Put("/", h.HandleUpdateUser)
				r.Delete("/", h.HandleDeleteUser)
			})
		})
	})
}

// HandleValidate Get User User.
// @Summary Gets a User User.
// @Description This endpoint Get User User.
// @Tags v1/User
// @Security JWTToken
// @Produce json
// @Success 200 {object} response.Base{data=user.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users [get]
func (h *UserHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	res, err := h.Service.GetByUserName(claims.UserName)

	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}

// HandleUpdateUser updates a User.
// @Summary updates a User.
// @Description This endpoint updates a User.
// @Tags v1/User
// @Security JWTToken
// @Param userId path string true "the user id"
// @Param User body user.NamePayload true "The Name to be changed"
// @Produce json
// @Success 200 {object} response.Base{data=user.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users/{userId} [put]
func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	userId, err := uuid.FromString(id)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var payload user.NamePayload
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

	res, err := h.Service.UpdateName(payload, userId)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}

// HandleDeleteUser Deletes a User.
// @Summary soft deletes a User.
// @Description This endpoint soft deletes a User.
// @Tags v1/User
// @Security JWTToken
// @Param userId path string true "the user id"
// @Produce json
// @Success 200 {object} response.Base{data=user.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/users/{userId} [delete]
func (h *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	userId, err := uuid.FromString(id)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	deleterId, err := uuid.FromString(claims.UserId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	res, err := h.Service.DeleteByID(userId, deleterId)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}
