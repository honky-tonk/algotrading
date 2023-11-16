package db

import (
	"algotrading/logger"
	"database/sql"

	_ "github.com/lib/pq"
)

func Db_main() *sql.DB {
	db, err := sql.Open("postgres", "user=algotrading_user password=000000 dbname=algotrading sslmode=disable")
	if err != nil {
		logger.Error.Fatal("can't open database: ", err)
	}
	return db
}
