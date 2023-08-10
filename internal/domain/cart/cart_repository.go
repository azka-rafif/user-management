package cart

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/jmoiron/sqlx"
)

type CartRepository interface {
	CartExistsByID(id string) (exists bool, err error)
	CreateItem(item CartItem) (err error)
	GetCartByID(cartId string) (res Cart, err error)
	GetCartItems(cartId string) (res []CartItem, err error)
	CartItemExistsByID(id string) (exists bool, err error)
	GetCartItemsByID(itemId string) (res CartItem, err error)
	Checkout(load []string) (err error)
	ProductExistsInCart(productId, cartId string) (exists bool, err error)
	GetCartItemByProduct(productId, cartId string) (res CartItem, err error)
	UpdateItem(item CartItem) (err error)
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
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txCreateItem(db, load); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *CartRepositoryMySQL) GetCartByID(cartId string) (res Cart, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM cart WHERE id = ?", cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) GetCartItems(cartId string) (res []CartItem, err error) {
	err = r.DB.Read.Select(&res, "SELECT * FROM cart_item WHERE cart_id = ?", cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) txCreateItem(tx *sqlx.Tx, load CartItem) (err error) {
	query := `INSERT INTO cart_item (id,cart_id,product_id,quantity,price,created_at,created_by,updated_at,updated_by)
	VALUES (:id,:cart_id,:product_id,:quantity,:price,:created_at,:created_by,:updated_at,:updated_by)`
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(load)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) Checkout(load []string) (err error) {
	if len(load) == 0 {
		return
	}
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		for _, item := range load {
			if err := r.txDeleteItem(db, item); err != nil {
				c <- err
				return
			}
		}
		c <- nil
	})
}

func (r *CartRepositoryMySQL) CartItemExistsByID(id string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(id) FROM cart_item WHERE id = ?", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) txDeleteItem(tx *sqlx.Tx, itemId string) (err error) {
	_, err = tx.Exec("DELETE FROM cart_item WHERE id = ?", itemId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) GetCartItemsByID(itemId string) (res CartItem, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM cart_item WHERE id = ?", itemId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) ProductExistsInCart(productId, cartId string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(product_id) FROM cart_item WHERE product_id = ? AND cart_id = ?", productId, cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) GetCartItemByProduct(productId, cartId string) (res CartItem, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM cart_item WHERE product_id = ? AND cart_id = ?", productId, cartId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *CartRepositoryMySQL) UpdateItem(item CartItem) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txUpdateItem(db, item); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *CartRepositoryMySQL) txUpdateItem(tx *sqlx.Tx, item CartItem) (err error) {

	query := `
	UPDATE cart_item
	SET
		quantity = :quantity,
		price = :price,
		created_at = :created_at,
		updated_at = :updated_at,
		deleted_at = :deleted_at,
		created_by = :created_by,
		updated_by = :updated_by,
		deleted_by = :deleted_by
	WHERE id = :id`
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		println(query)
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(item)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
