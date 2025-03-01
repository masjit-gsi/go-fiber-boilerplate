package database

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type DBConn struct {
	DB   *sqlx.DB
	Gorm *gorm.DB
}

func (d DBConn) Query() (db *sqlx.DB) {
	db = d.DB
	return
}

func (d DBConn) Orm() (db *gorm.DB) {
	db = d.Gorm
	return
}

// OpenDBConnection func for opening database connection.
func NewDBConnection() (dbConn DBConn, err error) {
	// Define Database connection variables.
	var (
		db   *sqlx.DB
		gorm *gorm.DB
	)

	// Get DB_TYPE value from .env file.
	dbType := os.Getenv("DB_TYPE")

	// Define a new Database connection with right DB type.
	switch dbType {
	case "pgx":
		db, err = PostgreSQLConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gorm, err = GormPostgreSQLConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	case "mysql":
		db, err = MysqlConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		gorm, err = GormMysqlConnection()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	dbConn = DBConn{
		DB:   db,
		Gorm: gorm,
	}

	return dbConn, nil
}
