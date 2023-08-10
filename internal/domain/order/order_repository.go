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
	GetAll(limit, offset int, sort, field, status string) (res []Order, err error)
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

func (r *OrderRepositoryMySQL) GetAll(limit, offset int, sort, field, status string) (res []Order, err error) {
	query := `SELECT * FROM atc_order `

	if status != "" {
		query += fmt.Sprintf("WHERE status = %s", status)
	}
	query += fmt.Sprintf("ORDER BY %s %s LIMIT %d OFFSET %d", field, sort, limit, offset)
	err = r.DB.Read.Select(&res, query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}
