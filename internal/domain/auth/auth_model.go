package auth

import (
	"encoding/json"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/roles"
)

type AuthPayload struct {
	Email    string `json:"email" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type LoginPayload struct {
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type JwtResponseFormat struct {
	AccessToken string `json:"access_token"`
}

func (j *JwtResponseFormat) MarshalJSON() ([]byte, error) {
	return json.Marshal(j)
}

func (p *AuthPayload) Validate() (err error) {
	validator := shared.GetValidator()
	p.Role = roles.GetStringFromRole(roles.GetRoleFromString(p.Role))
	return validator.Struct(p)
}
