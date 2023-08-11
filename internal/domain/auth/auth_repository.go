package auth

import (
	"github.com/evermos/boilerplate-go/infras"
)

type AuthRepository interface {
}

type AuthRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideAuthRepositoryMySQL(db *infras.MySQLConn) *AuthRepositoryMySQL {
	s := new(AuthRepositoryMySQL)
	s.DB = db
	return s
}
