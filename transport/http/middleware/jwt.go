package middleware

import (
	"context"
	"net/http"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type JwtAuthentication struct {
	conf *configs.Config
	db   *infras.MySQLConn
	jwt  *jwt.JWT
}

type ClaimsKey string

const (
	HeaderJwt = "Authorization"
)

func ProvideJwtAuthentication(conf *configs.Config, db *infras.MySQLConn) *JwtAuthentication {
	jwt := jwt.NewJWT(conf.App.JWTSecret)
	return &JwtAuthentication{
		conf: conf,
		db:   db,
		jwt:  jwt,
	}
}

func (a *JwtAuthentication) CheckJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderJwt)
		if token == "" {
			next.ServeHTTP(w, r)
		}
		_, err := a.jwt.ValidateJwt(token)
		if err != nil {
			next.ServeHTTP(w, r)
		}
		response.WithJSON(w, http.StatusOK, auth.JwtResponseFormat{AccessToken: token})
	})
}

func (a *JwtAuthentication) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(HeaderJwt)
		if token == "" {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		claims, err := a.jwt.ValidateJwt(token)
		if err != nil {
			response.WithError(w, failure.Unauthorized(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *JwtAuthentication) AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey("claims")).(*jwt.Claims)
		if !ok {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if claims.Role != "admin" {
			response.WithError(w, failure.Unauthorized("only admins are allowed"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *JwtAuthentication) CartAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "cartId")
		id, err := uuid.FromString(idString)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		claims, ok := r.Context().Value(ClaimsKey("claims")).(*jwt.Claims)
		if !ok {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if claims.CartId != id.String() && claims.Role != "admin" {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized, invalid credentials")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *JwtAuthentication) IsUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "userId")
		userId, err := uuid.FromString(id)
		if err != nil {
			response.WithError(w, failure.BadRequest(err))
			return
		}
		claims, ok := r.Context().Value(ClaimsKey("claims")).(*jwt.Claims)
		if !ok {
			response.WithMessage(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if userId.String() != claims.UserId && claims.Role != "admin" {
			response.WithError(w, failure.Unauthorized("Unauthorized, invalid credentials "))
			return
		}
		next.ServeHTTP(w, r)
	})
}
