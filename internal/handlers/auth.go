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
		r.Group(func(r chi.Router) {
			r.Post("/register", h.HandleRegister)
			r.Post("/login", h.HandleLogin)
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

func (h *AuthHandler) HandleValidate(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middleware.ClaimsKey("claims")).(*jwt.Claims)
	if !ok {
		response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	response.WithJSON(w, http.StatusOK, claims)
}

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
