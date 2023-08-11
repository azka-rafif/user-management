package user

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/evermos/boilerplate-go/internal/domain/cart"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/email"
	"github.com/evermos/boilerplate-go/shared/encrypt"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/evermos/boilerplate-go/shared/roles"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type User struct {
	UserId     uuid.UUID   `db:"id" validate:"required"`
	Email      string      `db:"email" validate:"required"`
	UserName   string      `db:"username" validate:"required"`
	Name       string      `db:"name" validate:"required"`
	Password   string      `db:"password" validate:"required"`
	Role       string      `db:"role" validate:"required"`
	CartId     uuid.UUID   `db:"cart_id" validate:"required"`
	Cart       cart.Cart   `db:"-"`
	Created_at time.Time   `db:"created_at" validate:"required"`
	Updated_at time.Time   `db:"updated_at" validate:"required"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type UserResponseFormat struct {
	UserId     uuid.UUID   `json:"id" validate:"required"`
	Email      string      `json:"email" validate:"required"`
	UserName   string      `json:"userName" validate:"required"`
	Name       string      `json:"name" validate:"required"`
	Password   string      `json:"-" validate:"required"`
	Role       string      `json:"role" validate:"required"`
	CartId     uuid.UUID   `json:"cartId" validate:"required"`
	Cart       cart.Cart   `json:"-"`
	Created_at time.Time   `json:"createdAt" validate:"required"`
	Updated_at time.Time   `json:"updatedAt" validate:"required"`
	Deleted_at null.Time   `json:"deletedAt"`
	Created_by uuid.UUID   `json:"createdBy"`
	Updated_by uuid.UUID   `json:"updatedBy"`
	Deleted_by nuuid.NUUID `json:"deletedBy"`
}

type UserPayload struct {
	Email    string `json:"email" validate:"required"`
	UserName string `json:"userName" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

func (u User) NewFromPayload(payload UserPayload) (res User, err error) {
	userId, err := uuid.NewV4()
	if err != nil {
		return
	}
	hashedPass, err := encrypt.HashPassword(payload.Password)
	if err != nil {
		return
	}
	userRole := roles.GetStringFromRole(roles.GetRoleFromString(payload.Role))
	cartId, err := uuid.NewV4()
	if err != nil {
		return
	}
	newCart, err := u.Cart.NewFromPayload(cart.CartPayload{CartId: cartId, UserId: userId})
	if err != nil {
		return
	}
	valid := email.Valid(payload.Email)
	if !valid {
		err = failure.BadRequest(errors.New("invalid email"))
		return
	}
	res = User{
		UserId:     userId,
		Email:      payload.Email,
		UserName:   payload.UserName,
		Name:       payload.Name,
		Password:   hashedPass,
		Role:       userRole,
		CartId:     cartId,
		Cart:       newCart,
		Created_at: time.Now().UTC(),
		Created_by: userId,
		Updated_at: time.Now().UTC(),
		Updated_by: userId,
	}
	err = res.Validate()
	return
}

type NamePayload struct {
	Name string `json:"name" validate:"required"`
}

func (u *User) ValidatePassword(loginPass string) error {
	return encrypt.ComparePasswords(u.Password, loginPass)
}

func (u *User) UpdateName(payload NamePayload) {
	u.Name = payload.Name
}

func (u User) ToResponseFormat() UserResponseFormat {
	return UserResponseFormat(u)
}

func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.ToResponseFormat())
}

func (u *User) Validate() (err error) {
	validator := shared.GetValidator()
	return validator.Struct(u)
}

func (u *User) isDeleted() bool {
	return u.Deleted_at.Valid && u.Deleted_by.Valid
}

func (u *User) SoftDelete(deleter uuid.UUID) (err error) {
	if u.isDeleted() {
		err = failure.Conflict("delete", "user", "already deleted")
		return
	}

	u.Deleted_at = null.TimeFrom(time.Now().UTC())
	u.Deleted_by = nuuid.From(deleter)

	err = u.Validate()

	return
}
