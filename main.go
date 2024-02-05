package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

const (
	dbDriver = "sqlite"
	dbUser   = "root"
	dbPass   = "toor"
	dbName   = "gocrud_app"
)

type User struct {
	ID    int
	Name  string
	Email string
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, "gocrud_app.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Parse JSON data from the request body
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	CreateUser(db, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Filed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")

}

func CreateUser(db *sql.DB, name, email string) error {
	query := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err := db.Exec(query, name, email)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	//Create a new router
	r := mux.NewRouter()

	// Define http routes using the router
	r.HandleFunc("/user", createUserHandler).Methods("POST")
	// r.HandleFunc("/user/{id}", getUserHandler).Methods("GET")
	// r.HandleFunc("/user/{id}", updateUserHandler).Methods("PUT")
	// r.HandleFunc("/user/{id}", deleteUserHandler).Methods("DELETE")

	//Start the http server on port 8090
	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
