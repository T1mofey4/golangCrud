package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, "gocrud_app.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// get id parametr from url
	vars := mux.Vars(r)
	idStr := vars["id"]

	// convert id to int
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		panic(err.Error())
	}

	user, err := GetUser(db, userID)
	if err != nil {
		http.Error(w, "User nor found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUser(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	row := db.QueryRow(query, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, "gocrud_app.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//Get the 'id' from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	//Convert 'id' to an integer
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)

	UpdateUser(db, userId, user.Name, user.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprintln(w, "User updated successfully")
}

func UpdateUser(db *sql.DB, id int, name, email string) error {
	query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
	fmt.Println(name, email)
	_, err := db.Exec(query, name, email, id)
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
	r.HandleFunc("/user/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/user/{id}", updateUserHandler).Methods("PUT")
	// r.HandleFunc("/user/{id}", deleteUserHandler).Methods("DELETE")

	//Start the http server on port 8090
	log.Println("Server listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", r))
}
