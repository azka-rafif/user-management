package product

import (
	"database/sql"
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	Create(prod Product) (err error)
	GetAll(limit, offset int, sort, field, productTitle string) (res []Product, err error)
	ExistsByID(id string) (exists bool, err error)
	GetByID(id string) (res Product, err error)
}

type ProductRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideProductRepositoryMySQL(db *infras.MySQLConn) *ProductRepositoryMySQL {
	return &ProductRepositoryMySQL{DB: db}
}

func (r *ProductRepositoryMySQL) Create(prod Product) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txCreate(db, prod); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *ProductRepositoryMySQL) txCreate(tx *sqlx.Tx, prod Product) (err error) {
	query := `INSERT INTO product (id,name,stock,price,created_at,created_by,updated_at,updated_by)
    VALUES (:id,:name,:stock,:price,:created_at,:created_by,:updated_at,:updated_by)`

	stmt, err := tx.PrepareNamed(query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(prod)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *ProductRepositoryMySQL) GetAll(limit, offset int, sort, field, productTitle string) (res []Product, err error) {
	query := `SELECT * FROM product `

	if productTitle != "" {
		query += `WHERE name `
		query += fmt.Sprintf("COLLATE UTF8_GENERAL_CI LIKE '%%%s%%' ", productTitle)
	}
	query += fmt.Sprintf("ORDER BY %s %s LIMIT %d OFFSET %d", field, sort, limit, offset)
	err = r.DB.Read.Select(&res, query)

	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return
}

func (r *ProductRepositoryMySQL) ExistsByID(id string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(id) FROM product WHERE id = ?", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *ProductRepositoryMySQL) GetByID(id string) (res Product, err error) {
	err = r.DB.Read.Get(&res, "SELECT * FROM product WHERE id = ?", id)
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("foo")
		logger.ErrorWithStack(err)
		return
	}
	return
}
