package order

import (
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
