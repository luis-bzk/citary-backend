package postgresql

import (
	"citary-backend/src/shared/constants"
	"database/sql"
	"log"
)

var DB *sql.DB

func Connect(database_url string) error {
	var err error
	DB, err = sql.Open("postgres", database_url)

	if err != nil {
		return err
	}

	// check connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// config pool connection
	DB.SetMaxOpenConns(constants.DBConstants.MaxOpenConnections)
	DB.SetMaxIdleConns(constants.DBConstants.MaxIdleConnections)

	log.Println("Connection to PostgreSQL established successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}

	return nil
}
