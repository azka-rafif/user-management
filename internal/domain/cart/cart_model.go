package cart

import (
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Cart struct {
	Id         uuid.UUID   `db:"id" validate:"required"`
	UserId     uuid.UUID   `db:"user_id" validate:"required"`
	Created_at time.Time   `db:"created_at" validate:"required"`
	Updated_at time.Time   `db:"updated_at" validate:"required"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type CartPayload struct {
	CartId uuid.UUID `json:"id" validate:"required"`
	UserId uuid.UUID `json:"userId" validate:"required"`
}

func (c Cart) NewFromPayload(user CartPayload) (res Cart, err error) {
	res = Cart{
		Id:         user.CartId,
		UserId:     user.UserId,
		Created_at: time.Now().UTC(),
		Created_by: user.UserId,
		Updated_at: time.Now().UTC(),
		Updated_by: user.UserId,
	}
	err = res.Validate()
	return
}

func (c *Cart) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(c)
}
