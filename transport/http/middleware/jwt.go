package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/internal/domain/auth"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/jwt"
	"github.com/evermos/boilerplate-go/transport/http/response"
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
		token = strings.Split(token, "Bearer ")[1]
		claims, err := a.jwt.ValidateJwt(token)
		if err != nil {
			response.WithError(w, failure.Unauthorized(err.Error()))
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
