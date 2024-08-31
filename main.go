package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	DAL "tmp/DAL"
	auth "tmp/auth"
	spel "tmp/speller"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Note struct {
	Content string `json:"content"`
}

type User struct {
	ID   string
	Name string
}

var Notes []Note

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "database"
)

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

func getNotes(w http.ResponseWriter, r *http.Request) {
	username, _ := auth.CheckToken(w, r)

	w.Header().Set("Content-Type", "application/json")

	db, _ := ConnectDB()
	defer db.Close()

	notes := []Note{}

	rows, err := db.Query(fmt.Sprintf("SELECT content FROM notes WHERE username = '%s'", username))
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

	json.NewEncoder(w).Encode(notes)
}

func deleteAllNotes(w http.ResponseWriter, r *http.Request) {
	username, _ := auth.CheckToken(w, r)

	w.Header().Set("Content-Type", "application/json")

	db, _ := ConnectDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM notes WHERE username = $1", username)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Все записи удалены.")
	}
	json.NewEncoder(w).Encode(Notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
	username, _ := auth.CheckToken(w, r)

	db, _ := ConnectDB()
	defer db.Close()

	var note Note
	_ = json.NewDecoder(r.Body).Decode(&note)
	Notes = append(Notes, note)
	json.NewEncoder(w).Encode(note)
	fmt.Println(note.Content)
	checkContent, _ := spel.CheckText(note.Content)

	_, err := db.Exec(fmt.Sprintf("INSERT INTO notes (content,username) VALUES ('%s','%s')", checkContent, username))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Запись успешно добавлена!\n")

	}

}

// Авторизация пользователя

func handleRequest(w http.ResponseWriter, r *http.Request) {

	type Data struct {
		Username string `json:"username"`
	}

	username := r.PostFormValue("username")
	var data Data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Fprintf(w, "Ошибка при декодировании запроса: %v", err)
		return
	}
	fmt.Println(data.Username)
	if DAL.CheckUser(data.Username) == true {
		fmt.Println("Пользователь найден")
		fmt.Println(username)
		token, err := auth.GenerateToken(data.Username)
		fmt.Println(token)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(24 * time.Hour),
		})

	} else {

		fmt.Println("error")
	}

}

func main() {
	db, _ := ConnectDB()
	defer db.Close()

	r := mux.NewRouter()
	log.Println("Listening at port 8082")
	r.HandleFunc("/notes", getNotes).Methods("GET")
	r.HandleFunc("/notes", deleteAllNotes).Methods("DELETE")
	r.HandleFunc("/notes", createNote).Methods("POST")
	r.HandleFunc("/login", handleRequest).Methods("POST")
	log.Fatal(http.ListenAndServe(":8082", r))

}
