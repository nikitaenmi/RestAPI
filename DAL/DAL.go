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

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, err
}

func CheckUser(username string) bool {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()

	notes := []Note{}

	rows, err := db.Query(fmt.Sprintf("SELECT id FROM users WHERE users = '%s'", username))
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		p := Note{}
		err := rows.Scan(&p.Content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		notes = append(notes, p)

	}
	if len(notes) < 1 {
		fmt.Println("Не нашлость")
		fmt.Println(notes)
		return false
	}
	if len(notes) >= 1 {
		fmt.Println("Нашлось!")
		return true
	}
	fmt.Println(notes)
	return true
}
