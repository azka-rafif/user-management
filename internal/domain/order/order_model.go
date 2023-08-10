package order

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Order struct {
	Id         uuid.UUID   `db:"id" validate:"required"`
	UserId     uuid.UUID   `db:"user_id" validate:"required"`
	TotalPrice float64     `db:"total_price" validate:"required"`
	Status     string      `db:"status" validate:"required"`
	OrderItems []OrderItem `db:"-"`
	CreatedAt  time.Time   `db:"created_at" validate:"required"`
	UpdatedAt  time.Time   `db:"updated_at" validate:"required"`
	DeletedAt  null.Time   `db:"deleted_at"`
	CreatedBy  uuid.UUID   `db:"created_by" validate:"required"`
	UpdatedBy  uuid.UUID   `db:"updated_by" validate:"required"`
	DeletedBy  nuuid.NUUID `db:"deleted_by"`
}

type OrderItem struct {
	Id        uuid.UUID   `db:"id" validate:"required"`
	OrderId   uuid.UUID   `db:"order_id" validate:"required"`
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

type OrderResponseFormat struct {
	Id         uuid.UUID   `json:"id" validate:"required"`
	UserId     uuid.UUID   `json:"userId" validate:"required"`
	TotalPrice float64     `json:"totalPrice" validate:"required"`
	Status     string      `json:"status" validate:"required"`
	OrderItems []OrderItem `json:"-"`
	CreatedAt  time.Time   `json:"createdAt" validate:"required"`
	UpdatedAt  time.Time   `json:"updatedAt" validate:"required"`
	DeletedAt  null.Time   `json:"deletedAt"`
	CreatedBy  uuid.UUID   `json:"createdBy" validate:"required"`
	UpdatedBy  uuid.UUID   `json:"updatedBy" validate:"required"`
	DeletedBy  nuuid.NUUID `json:"deletedBy"`
}

type OrderItemResponseFormat struct {
	Id        uuid.UUID   `json:"id" validate:"required"`
	OrderId   uuid.UUID   `json:"orderId" validate:"required"`
	ProductId uuid.UUID   `json:"productId" validate:"required"`
	Quantity  int         `json:"quantity" validate:"required"`
	Price     float64     `json:"price" validate:"required"`
	CreatedAt time.Time   `json:"createdAt" validate:"required"`
	UpdatedAt time.Time   `json:"updatedAt" validate:"required"`
	DeletedAt null.Time   `json:"deletedAt"`
	CreatedBy uuid.UUID   `json:"createdBy" validate:"required"`
	UpdatedBy uuid.UUID   `json:"updatedBy" validate:"required"`
	DeletedBy nuuid.NUUID `json:"deletedBy"`
}

type OrderPayload struct {
	UserId     uuid.UUID
	TotalPrice float64
	Status     string
}

type OrderItemPayload struct {
	ProductId uuid.UUID
	UserId    uuid.UUID
	OrderId   uuid.UUID
	Quantity  int
	Price     float64
}

func (o Order) NewFromPayload(load OrderPayload) (res Order, err error) {
	orderId, err := uuid.NewV4()
	if err != nil {
		return
	}
	res = Order{
		Id:         orderId,
		UserId:     load.UserId,
		TotalPrice: load.TotalPrice,
		Status:     load.Status,
		CreatedAt:  time.Now().UTC(),
		CreatedBy:  load.UserId,
		UpdatedAt:  time.Now().UTC(),
		UpdatedBy:  load.UserId,
	}
	err = res.Validate()
	return
}

func (o OrderItem) NewFromPayload(load OrderItemPayload) (res OrderItem, err error) {
	orderItemId, err := uuid.NewV4()
	if err != nil {
		return
	}
	res = OrderItem{
		Id:        orderItemId,
		OrderId:   load.OrderId,
		ProductId: load.ProductId,
		Quantity:  load.Quantity,
		Price:     load.Price,
		CreatedAt: time.Now().UTC(),
		CreatedBy: load.UserId,
		UpdatedAt: time.Now().UTC(),
		UpdatedBy: load.UserId,
	}
	err = res.Validate()
	return
}

func (o *Order) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(o)
}

func (o *OrderItem) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(o)
}

func (o *Order) AttachItems(load []OrderItem) Order {
	o.OrderItems = append(o.OrderItems, load...)
	return *o
}

func (o Order) ToResponseFormat() OrderResponseFormat {
	return OrderResponseFormat(o)
}

func (o OrderItem) ToResponseFormat() OrderItemResponseFormat {
	return OrderItemResponseFormat(o)
}

func (o *Order) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ToResponseFormat())
}

func (o *OrderItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ToResponseFormat())
}
