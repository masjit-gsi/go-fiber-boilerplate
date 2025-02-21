package repository

import (
	"github.com/fiber-go-template/app/models"
	"github.com/fiber-go-template/config/logger"
	"github.com/fiber-go-template/database"
)

var (
	userQuery = struct {
		Select string
	}{
		Select: `SELECT * FROM users `,
	}
)

type UserRepositoryDB struct {
	DB database.DBConn
}

func NewUserRepository(db database.DBConn) UserRepository {
	return &UserRepositoryDB{
		DB: db,
	}
}

type UserRepository interface {
	GetUserByID(id string) (user models.User, err error)
	GetUserByUsername(username string) (user models.User, err error)
}

// GetUserByID query for getting one User by given ID.
func (r *UserRepositoryDB) GetUserByID(id string) (user models.User, err error) {
	err = r.DB.Query().Get(&user, userQuery.Select+" where id=$1", id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return user, nil
}

// GetUserByEmail query for getting one User by given Username.
func (r *UserRepositoryDB) GetUserByUsername(username string) (user models.User, err error) {
	err = r.DB.Query().Get(&user, userQuery.Select+" where email=$1", username)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	return user, nil
}
