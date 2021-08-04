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

var conn *pgx.Conn // global variable connection so we wont make a new connection for each query. (good for small projects)

const (
	host     = "localhost" // default
	port     = 5432        // default
	user     = "postgres"  // default
	password = "Edit your password"
	db_name  = "Edit your DataBase Name"
)

type User struct {
	Name      string `json:"name"`
	LastName  string `json:"lastname"`
	BirthDate string `json:"birthdate"`
}

type JsonResponse struct {
	Type    string `json:"type"`
	Data    []User `json:"data"`
	Message string `json:"message"`
}

func main() {

	router := mux.NewRouter()
	var err error
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, db_name)
	conn, err = pgx.Connect(context.Background(), dbinfo) // connect to psql
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
	fmt.Println("listening on 8080 port")

	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users", GetUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Accept", "application/json")

	var My_user User
	err := json.NewDecoder(r.Body).Decode(&My_user)
	checkErr(err)

	var response = JsonResponse{}

	if My_user.Name == "" || My_user.LastName == "" || My_user.BirthDate == "" {
		response = JsonResponse{Type: "error", Data: []User{My_user}, Message: "parameters are missing from Post request!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var lastInsertID int
	err = conn.QueryRow(context.Background(), "INSERT INTO users(name, lastname, birthdate) VALUES($1, $2, $3) returning id;", My_user.Name, My_user.LastName, My_user.BirthDate).Scan(&lastInsertID)
	checkErr(err)

	response = JsonResponse{Type: "success", Data: []User{My_user}, Message: "The user has been inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
