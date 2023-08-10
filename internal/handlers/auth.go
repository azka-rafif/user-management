package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
)

type AuthHandler struct {
	Service auth.AuthService
	JwtAuth *middleware.JwtAuthentication
}

func ProvideAuthHandler(service auth.AuthService, jwtAuth *middleware.JwtAuthentication) AuthHandler {
	return AuthHandler{Service: service, JwtAuth: jwtAuth}
}

func (h *AuthHandler) Router(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.HandleLogin)
		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.Validate)
			r.Use(h.JwtAuth.AdminOnly)
			r.Post("/register", h.HandleRegister)
		})
		r.Group(func(r chi.Router) {
			r.Use(h.JwtAuth.Validate)
			r.Get("/validate", h.HandleValidate)
		})
	})
	r.Route("/profile", func(r chi.Router) {
		r.Use(h.JwtAuth.Validate)
		r.Get("/", h.HandleGetProfile)
		r.Put("/", h.HandleUpdateProfile)
	})
}

// HandleRegister creates a new User.
// @Summary Create a new User / register a user.
// @Description This endpoint creates a new User.
// @Tags v1/Auth
// @Security JWTToken
// @Param User body auth.AuthPayload true "The User to be created."
// @Produce json
// @Success 201 {object} response.Base{data=auth.JwtResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/auth/register [post]
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload auth.AuthPayload
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

	res, err := h.Service.Register(payload)
	if err != nil {
		response.WithError(w, err)
		return
	}
	response.WithJSON(w, http.StatusCreated, res)
}

// HandleLogin Login a user.
// @Summary Login a user.
// @Description This endpoint Logs in a User.
// @Tags v1/Auth
// @Param User body auth.LoginPayload true "The User to be logged in."
// @Produce json
// @Success 201 {object} response.Base{data=auth.JwtResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/auth/login [post]
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload auth.LoginPayload
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

	res, err := h.Service.Login(payload)
	if err != nil {
		response.WithError(w, failure.InternalError(err))
		return
	}

	response.WithJSON(w, http.StatusOK, res)
}

// HandleValidate validates a JWT Token.
// @Summary Validates the given Jwt Token.
// @Description This endpoint validates a jwt token.
// @Tags v1/Auth
// @Security JWTToken
// @Produce json
// @Success 200 {object} response.Base{data=auth.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/auth/validate [get]
func (h *AuthHandler) HandleValidate(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response.WithJSON(w, http.StatusOK, claims)
}

// HandleValidate Get User Profile.
// @Summary Gets a User Profile.
// @Description This endpoint Get User Profile.
// @Tags v1/Profile
// @Security JWTToken
// @Produce json
// @Success 200 {object} response.Base{data=auth.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/profile [get]
func (h *AuthHandler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
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

// HandleUpdateProfile updates a users profile.
// @Summary updates a users profile.
// @Description This endpoint updates a users profile.
// @Tags v1/Profile
// @Security JWTToken
// @Param User body auth.NamePayload true "The Name to be changed"
// @Produce json
// @Success 209 {object} response.Base{data=auth.UserResponseFormat}
// @Failure 400 {object} response.Base
// @Failure 409 {object} response.Base
// @Failure 500 {object} response.Base
// @Router /v1/profile [put]
func (h *AuthHandler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload auth.NamePayload
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
	res, err := h.Service.UpdateName(payload, claims.UserName)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}
	response.WithJSON(w, http.StatusOK, res)
}
