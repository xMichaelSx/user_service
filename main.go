package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your_password"
	db_name  = "test"
)

type User struct {
	UserID    string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastname"`
	Age       string `json:"age"`
	BirthDate string `json:"birthdate"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

func main() {

	router := mux.NewRouter()

	// Create a user
	router.HandleFunc("/users", CreateUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":5432", router))

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting...")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	user.UserID = r.Form.Get("id")
	user.Name = r.Form.Get("name")
	user.LastName = r.Form.Get("lastname")
	user.Age = r.Form.Get("age")
	user.BirthDate = r.Form.Get("birthdate")

	var response = JsonResponse{}

	if user.UserID == "" || user.Name == "" || user.LastName == "" || user.Age == "" || user.BirthDate == "" {
		response = JsonResponse{Type: "error", Message: "parameters are missing!"}
	} else {

		dbinfo := fmt.Sprintf("postgresql://host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db_name)
		conn, err := pgx.Connect(context.Background(), dbinfo) // connect to psql
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close(context.Background())

		var lastInsertID int
		err = conn.QueryRow(context.Background(), "INSERT INTO users(id, name, last, age, birthdate) VALUES($1, $2, $3, $4, $5) returning id;", user.UserID, user.Name, user.LastName, user.Age, user.BirthDate).Scan(&lastInsertID)
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The user has been inserted successfully!"}

	}

	json.NewEncoder(w).Encode(response)
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
