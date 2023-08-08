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
	CartItems  []CartItem  `db:"-"`
	Created_at time.Time   `db:"created_at" validate:"required"`
	Updated_at time.Time   `db:"updated_at" validate:"required"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type CartItem struct {
	Id         uuid.UUID   `db:"id" validate:"required"`
	CartId     uuid.UUID   `db:"cart_id" validate:"required"`
	ProductId  uuid.UUID   `db:"product_id" validate:"required"`
	Quantity   int         `db:"quantity" validate:"required"`
	Price      float64     `db:"price" validate:"required"`
	Created_at time.Time   `db:"created_at" validate:"required"`
	Updated_at time.Time   `db:"updated_at" validate:"required"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type CartItemPayload struct {
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}

type CartPayload struct {
	CartId uuid.UUID `json:"id" validate:"required"`
	UserId uuid.UUID `json:"userId" validate:"required"`
}

func (c Cart) NewFromPayload(load CartPayload) (res Cart, err error) {
	res = Cart{
		Id:         load.CartId,
		UserId:     load.UserId,
		Created_at: time.Now().UTC(),
		Created_by: load.UserId,
		Updated_at: time.Now().UTC(),
		Updated_by: load.UserId,
	}
	err = res.Validate()
	return
}

func (c CartItem) NewFromPayload(load CartItemPayload, cartId, userId uuid.UUID, productPrice float64) (res CartItem, err error) {
	cartItemId, err := uuid.NewV4()
	if err != nil {
		return
	}
	res = CartItem{
		Id:         cartItemId,
		CartId:     cartId,
		ProductId:  load.ProductId,
		Quantity:   load.Quantity,
		Price:      productPrice * float64(load.Quantity),
		Created_at: time.Now().UTC(),
		Created_by: userId,
		Updated_at: time.Now().UTC(),
		Updated_by: userId,
	}
	err = res.Validate()
	return
}

func (c *Cart) AttachItems(load []CartItem) Cart {
	c.CartItems = append(c.CartItems, load...)
	return *c
}
func (c *CartItem) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(c)
}

func (c *Cart) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(c)
}
