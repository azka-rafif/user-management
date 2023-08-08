package cart

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
)

type CartRepository interface {
	CartExistsByID(id string) (exists bool, err error)
	CreateItem(item CartItem) (err error)
	GetCartByID(cartId string) (res Cart, err error)
}

type CartRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideCartRepositoryMySQL(db *infras.MySQLConn) *CartRepositoryMySQL {
	return &CartRepositoryMySQL{DB: db}
}

func (r *CartRepositoryMySQL) CartExistsByID(id string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(id) FROM cart WHERE id = ?", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) CreateItem(load CartItem) (err error) {
	return
}

func (r *CartRepositoryMySQL) GetCartByID(cartId string) (res Cart, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM cart WHERE id = ?", cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
