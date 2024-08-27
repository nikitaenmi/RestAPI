package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Note struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	UserID  int
}

var Notes []Note

func getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Notes)
}

func deleteAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.Ping()
	CheckError(err)

	_, err = db.Exec("DELETE FROM notes;")
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Все записи удалены.")
	}
	json.NewEncoder(w).Encode(Notes)
}

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "database"
)

func main() {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.Ping()
	CheckError(err)

	r := mux.NewRouter()
	log.Println("Listening at port 8080")
	r.HandleFunc("/notes", getNotes).Methods("GET")
	r.HandleFunc("/", deleteAllNotes).Methods("DELETE")
	r.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":8080", r))

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./public/html/startStaticPage.html")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()
	err = db.Ping()

	if r.Method != http.MethodPost {
		t.Execute(w, nil) // отображает страницу
		return
	}

	r.ParseForm()
	text := r.PostFormValue("text") // здесь запрос
	CheckError(err)
	fmt.Println("Connected!")
	fmt.Println(text)

	_, err = db.Exec(fmt.Sprintf("INSERT INTO notes (userid, content) VALUES (1, '%s')", text))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Запись успешно добавлена!\n")
	}

	rows, _ := db.Query("SELECT content FROM notes WHERE userID = 1")
	for rows.Next() {
		p := Note{}
		err := rows.Scan(&p.Content)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Notes = append(Notes, p)
	}

	tmpl, _ := template.ParseFiles("./public/html/startStaticPage.html")
	tmpl.Execute(w, Notes)

}
