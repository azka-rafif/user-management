package product

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/nuuid"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Product struct {
	Id         uuid.UUID   `db:"id" validate:"required"`
	Name       string      `db:"name" validate:"required"`
	Stock      int         `db:"stock" validate:"required"`
	Price      float64     `db:"price" validate:"required"`
	Created_at time.Time   `db:"created_at" validate:"required"`
	Updated_at time.Time   `db:"updated_at" validate:"required"`
	Deleted_at null.Time   `db:"deleted_at"`
	Created_by uuid.UUID   `db:"created_by"`
	Updated_by uuid.UUID   `db:"updated_by"`
	Deleted_by nuuid.NUUID `db:"deleted_by"`
}

type ProductResponseFormat struct {
	Id         uuid.UUID   `json:"id" validate:"true"`
	Name       string      `json:"name" validate:"true"`
	Stock      int         `json:"stock" validate:"true"`
	Price      float64     `json:"price" validate:"true"`
	Created_at time.Time   `json:"createdAt" validate:"required"`
	Updated_at time.Time   `json:"updatedAt" validate:"required"`
	Deleted_at null.Time   `json:"deletedAt,omitempty"`
	Created_by uuid.UUID   `json:"createdBy"`
	Updated_by uuid.UUID   `json:"updatedBy"`
	Deleted_by nuuid.NUUID `json:"deletedBy,omitempty"`
}

type ProductPayload struct {
	Name  string  `json:"name"`
	Stock int     `json:"stock"`
	Price float64 `json:"price"`
}

func (p Product) NewFromPayload(load ProductPayload, userId uuid.UUID) (res Product, err error) {

	prodId, err := uuid.NewV4()

	if err != nil {
		return
	}

	res = Product{
		Id:         prodId,
		Name:       load.Name,
		Stock:      load.Stock,
		Price:      load.Price,
		Created_at: time.Now().UTC(),
		Created_by: userId,
		Updated_at: time.Now().UTC(),
		Updated_by: userId,
	}
	err = res.Validate()
	return
}

func (p *Product) Validate() error {
	validator := shared.GetValidator()
	return validator.Struct(p)
}

func (p Product) ToResponseFormat() ProductResponseFormat {
	return ProductResponseFormat(p)
}

func (p *Product) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponseFormat())
}
