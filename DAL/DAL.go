package dal

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "database"
)

type Note struct {
	ID       int64
	Username string
	Content  string
}

func ConnectDB() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return db, err
}

func CreateTable() {
	db, err := ConnectDB()
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS notes (id SERIAL PRIMARY KEY, content TEXT, username TEXT)")
	if err != nil {
		log.Fatal(err)
	}

}
