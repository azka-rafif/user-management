package auth

import (
	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	Create(user User) (err error)
	ExistsByID(userId uuid.UUID) (exists bool, err error)
	ExistsByUserName(userName string) (exists bool, err error)
	GetByUserName(userName string) (user User, err error)
	Update(user User) (err error)
}

type AuthRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideAuthRepositoryMySQL(db *infras.MySQLConn) *AuthRepositoryMySQL {
	s := new(AuthRepositoryMySQL)
	s.DB = db
	return s
}

func (r *AuthRepositoryMySQL) Create(user User) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txCreate(db, user); err != nil {
			c <- err
			return
		}
		if err := r.txCreateCart(db, user); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *AuthRepositoryMySQL) ExistsByID(userId uuid.UUID) (exists bool, err error) {

	err = r.DB.Read.Get(&exists, "SELECT COUNT(id) FROM user WHERE id = ?", userId.String())

	if err != nil {
		err = failure.NotFound("user")
		logger.ErrorWithStack(err)
	}

	return
}

func (r *AuthRepositoryMySQL) ExistsByUserName(userName string) (exists bool, err error) {
	err = r.DB.Read.Get(&exists, "SELECT COUNT(username) FROM user WHERE username = ?", userName)
	if err != nil {
		err = failure.NotFound("user")
		logger.ErrorWithStack(err)
	}

	return
}

func (r *AuthRepositoryMySQL) GetByUserName(userName string) (user User, err error) {
	exists, err := r.ExistsByUserName(userName)
	if err != nil {
		return
	}
	if !exists {
		err = failure.NotFound("user")
		return
	}
	err = r.DB.Read.Get(&user, "SELECT * from user WHERE username = ?", userName)
	if err != nil {
		err = failure.NotFound("user")
		return
	}
	return
}

func (r *AuthRepositoryMySQL) txCreate(tx *sqlx.Tx, payload User) (err error) {
	query := `insert into user (id,email,username,name,password,role,cart_id,created_at,created_by,updated_at,updated_by)
    VALUES (:id,:email,:username,:name,:password,:role,:cart_id,:created_at,:created_by,:updated_at,:updated_by)`

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(payload)
	if err != nil {
		// err = failure.Conflict("create", "user", "already exist")
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *AuthRepositoryMySQL) txCreateCart(tx *sqlx.Tx, payload User) (err error) {
	query := `INSERT INTO cart (id,user_id,created_at,created_by,updated_at,updated_by)
	VALUES (:id,:user_id,:created_at,:created_by,:updated_at,:updated_by)`
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(payload.Cart)
	if err != nil {
		err = failure.Conflict("create", "cart", "already exist")
		logger.ErrorWithStack(err)
		return
	}
	return
}

func (r *AuthRepositoryMySQL) Update(user User) (err error) {
	return r.DB.WithTransaction(func(db *sqlx.Tx, c chan error) {
		if err := r.txUpdate(db, user); err != nil {
			c <- err
			return
		}
		c <- nil
	})
}

func (r *AuthRepositoryMySQL) txUpdate(tx *sqlx.Tx, payload User) (err error) {
	query := `UPDATE user
	SET 
		id = :id,
		username = :username,
		name = :name,
		password = :password,
		role =  :role,
		created_at = :created_at,
		created_by = :created_by,
		updated_at = :updated_at,
		updated_by = :updated_by,
		deleted_at = :deleted_at,
		deleted_by = :deleted_by
	WHERE id = :id`
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(payload)
	if err != nil {
		tx.Rollback()
		logger.ErrorWithStack(err)
		return
	}
	return
}
