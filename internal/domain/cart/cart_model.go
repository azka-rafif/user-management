package cart

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Cart struct {
	Id        uuid.UUID   `db:"id" validate:"required"`
	UserId    uuid.UUID   `db:"user_id" validate:"required"`
	CartItems []CartItem  `db:"-"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedBy uuid.UUID   `db:"updated_by" validate:"required"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

type CartItem struct {
	Id        uuid.UUID   `db:"id" validate:"required"`
	CartId    uuid.UUID   `db:"cart_id" validate:"required"`
	ProductId uuid.UUID   `db:"product_id" validate:"required"`
	Quantity  int         `db:"quantity" validate:"required"`
	Price     float64     `db:"price" validate:"required"`
	CreatedAt time.Time   `db:"created_at" validate:"required"`
	UpdatedAt time.Time   `db:"updated_at" validate:"required"`
	DeletedAt null.Time   `db:"deleted_at"`
	CreatedBy uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedBy uuid.UUID   `db:"updated_by" validate:"required"`
	DeletedBy nuuid.NUUID `db:"deleted_by"`
}

type CartItemPayload struct {
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
}

type CartPayload struct {
	CartId uuid.UUID `json:"id" validate:"required"`
	UserId uuid.UUID `json:"userId" validate:"required"`
}

type CheckoutPayload struct {
	CartItemsIds []string `json:"items" validate:"required"`
}

type CartResponseFormat struct {
	Id        uuid.UUID                `json:"id" validate:"required"`
	UserId    uuid.UUID                `json:"userId" validate:"required"`
	CartItems []CartItemResponseFormat `json:"cartItems,omitempty"`
	CreatedAt time.Time                `json:"createdAt" validate:"required"`
	UpdatedAt time.Time                `json:"updatedAt" validate:"required"`
	DeletedAt null.Time                `json:"deletedAt,omitempty"`
	CreatedBy uuid.UUID                `json:"createdBy"`
	UpdatedBy uuid.UUID                `json:"updatedBy"`
	DeletedBy nuuid.NUUID              `json:"deletedBy,omitempty"`
}

type CartItemResponseFormat struct {
	Id        uuid.UUID   `json:"id" validate:"required"`
	CartId    uuid.UUID   `json:"cartId" validate:"required"`
	ProductId uuid.UUID   `json:"productId" validate:"required"`
	Quantity  int         `json:"quantity" validate:"required"`
	Price     float64     `json:"price" validate:"required"`
	CreatedAt time.Time   `json:"createdAt" validate:"required"`
	UpdatedAt time.Time   `json:"updatedAt" validate:"required"`
	DeletedAt null.Time   `json:"deletedAt,omitempty"`
	CreatedBy uuid.UUID   `json:"createdBy"`
	UpdatedBy uuid.UUID   `json:"updatedBy"`
	DeletedBy nuuid.NUUID `json:"deletedBy,omitempty"`
}

func (c Cart) NewFromPayload(load CartPayload) (res Cart, err error) {
	res = Cart{
		Id:        load.CartId,
		UserId:    load.UserId,
		CreatedAt: time.Now().UTC(),
		CreatedBy: load.UserId,
		UpdatedAt: time.Now().UTC(),
		UpdatedBy: load.UserId,
	}
	err = res.Validate()
	return
}

func (c CartItem) NewFromPayload(load CartItemPayload, cartId, userId uuid.UUID, productPrice float64) (res CartItem, err error) {
	cartItemId, err := uuid.NewV4()
	if err != nil {
		return
	}
	cartItemPrice := productPrice * float64(load.Quantity)
	res = CartItem{
		Id:        cartItemId,
		CartId:    cartId,
		ProductId: load.ProductId,
		Quantity:  load.Quantity,
		Price:     cartItemPrice,
		CreatedAt: time.Now().UTC(),
		CreatedBy: userId,
		UpdatedAt: time.Now().UTC(),
		UpdatedBy: userId,
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

func (c *Cart) ToResponseFormat() CartResponseFormat {
	var items []CartItemResponseFormat
	for _, item := range c.CartItems {
		items = append(items, item.ToResponseFormat())
	}
	res := CartResponseFormat{
		Id:        c.Id,
		UserId:    c.UserId,
		CartItems: items,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
		CreatedBy: c.CreatedBy,
		UpdatedBy: c.UpdatedBy,
		DeletedBy: c.DeletedBy,
	}

	return res
}

func (c CartItem) ToResponseFormat() (res CartItemResponseFormat) {
	return CartItemResponseFormat(c)
}

func (c Cart) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.ToResponseFormat())
}
