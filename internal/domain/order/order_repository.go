package order

import (
	"fmt"
	"strings"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/jmoiron/sqlx"
)

type OrderRepository interface {
	Create(load Order) (err error)
	GetAll(limit, offset int, sort, field, status, userId, userRole string, cancelled bool) (res []Order, err error)
	CancelOrder(load Order) (err error)
	GetOrderByID(orderId string) (res Order, err error)
	ExistsByID(orderId string) (exists bool, err error)
	GetItemsByOrderID(orderId string) (res []OrderItem, err error)
}

type OrderRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideOrderRepositoryMySQL(db *infras.MySQLConn) *OrderRepositoryMySQL {
	return &OrderRepositoryMySQL{DB: db}
}

func (r *OrderRepositoryMySQL) Create(load Order) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txCreate(db, load); err != nil {
			c <- err
			return
		}
		if err := r.txCreateItems(db, load.OrderItems); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *OrderRepositoryMySQL) txCreate(tx *sqlx.Tx, load Order) (err error) {
	query := `INSERT INTO atc_order (id,user_id,total_price,status,created_at,updated_at,created_by,updated_by) 
	VALUES (:id,:user_id,:total_price,:status,:created_at,:updated_at,:created_by,:updated_by)`
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

func (r *OrderRepositoryMySQL) txCreateItems(tx *sqlx.Tx, load []OrderItem) (err error) {
	if len(load) == 0 {
		return
	}
	query, args, err := r.composeBulkInsertItemQuery(load)
	if err != nil {
		return
	}
	stmt, err := tx.Preparex(query)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *OrderRepositoryMySQL) composeBulkInsertItemQuery(load []OrderItem) (query string, params []interface{}, err error) {
	values := []string{}
	for _, oi := range load {
		param := map[string]interface{}{
			"id":         oi.Id,
			"order_id":   oi.OrderId,
			"product_id": oi.ProductId,
			"quantity":   oi.Quantity,
			"price":      oi.Price,
			"created_at": oi.CreatedAt,
			"updated_at": oi.UpdatedAt,
			"created_by": oi.CreatedBy,
			"updated_by": oi.UpdatedBy,
		}
		q, args, err := sqlx.Named(`(:id,:order_id,:product_id,:quantity,:price,:created_at,:updated_at,:created_by,:updated_by)`, param)
		if err != nil {
			return query, params, err
		}
		values = append(values, q)
		params = append(params, args...)
	}
	query = fmt.Sprintf(`%v %v`, `INSERT INTO order_item (
				id,
				order_id,
				product_id,
				quantity,
				price,
				created_at,
				updated_at,
				created_by,
				updated_by
			) VALUES `, strings.Join(values, ","))
	return
}

func (r *OrderRepositoryMySQL) GetAll(limit, offset int, sort, field, status, userId, userRole string, cancelled bool) (res []Order, err error) {
	query := `SELECT * FROM atc_order `

	if userRole != "admin" {
		query += fmt.Sprintf("WHERE user_id = '%s' ", userId)
	}
	if status != "" {
		exists := strings.Contains(query, "WHERE")
		if !exists && userRole != "admin" {
			query += "WHERE "
		}
		if exists && strings.Contains(query, "user_id") {
			query += "AND "
		}
		query += fmt.Sprintf("status = '%s' ", status)
	}
	if cancelled {
		exists := strings.Contains(query, "WHERE")
		if !exists {
			query += `WHERE deleted_at IS NOT NULL `
		}
		if exists {
			query += `AND deleted_at IS NOT NULL `
		}

	}
	query += fmt.Sprintf("ORDER BY %s %s LIMIT %d OFFSET %d", field, sort, limit, offset)
	err = r.DB.Read.Select(&res, query)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *OrderRepositoryMySQL) CancelOrder(load Order) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txUpdate(db, load); err != nil {
			c <- err
			return
		}
		for _, item := range load.OrderItems {
			if err := r.txUpdateItem(db, item); err != nil {
				c <- err
				return
			}
		}
		c <- nil
	})
}

func (r *OrderRepositoryMySQL) txUpdate(tx *sqlx.Tx, load Order) (err error) {
	query := `
	UPDATE atc_order
	SET
		total_price = :total_price,
		status = :status,
		created_at = :created_at,
		updated_at = :updated_at,
		deleted_at = :deleted_at,
		created_by = :created_by,
		updated_by = :updated_by,
		deleted_by = :deleted_by
	WHERE id = :id`
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

func (r *OrderRepositoryMySQL) txUpdateItem(tx *sqlx.Tx, load OrderItem) (err error) {
	query := `
	UPDATE order_item
	SET
		price = :price,
		quantity = :quantity,
		created_at = :created_at,
		updated_at = :updated_at,
		deleted_at = :deleted_at,
		created_by = :created_by,
		updated_by = :updated_by,
		deleted_by = :deleted_by
	WHERE id = :id`
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

func (r *OrderRepositoryMySQL) GetOrderByID(orderId string) (res Order, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM atc_order WHERE id = ?", orderId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *OrderRepositoryMySQL) ExistsByID(orderId string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(id) FROM atc_order WHERE id = ?", orderId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *OrderRepositoryMySQL) GetItemsByOrderID(orderId string) (res []OrderItem, err error) {
	err = r.DB.Read.Select(&res, "SELECT * FROM order_item WHERE order_id = ?", orderId)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
