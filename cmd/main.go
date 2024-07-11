package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

const (
	dbDriver = "mysql"
	dbUser   = "root"
	dbPass   = "root"
	dbName   = "db_golangcrud"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/user", createUserHandle).Methods("POST")
	router.HandleFunc("/user", getAllUsersHandle).Methods("GET")
	router.HandleFunc("/user/{id}", getUserHandle).Methods("GET")
	router.HandleFunc("/user/{id}", updateUserHandle).Methods("PUT")
	router.HandleFunc("/user/{id}", deleteUserHandle).Methods("DELETE")

	log.Println("Starting server on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

type User struct {
	ID    int
	Name  string
	Email string
}

func createUserHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	CreateUser(db, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User created successfully")
}

func CreateUser(db *sql.DB, name string, email string) error {
	createUserQuery := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err := db.Exec(createUserQuery, name, email)
	if err != nil {
		return err
	}
	return nil
}

func getAllUsersHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	user, err := GetAllUsers(db)
	if err != nil {
		http.Error(w, "Error gettin users", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	getAllUsersQuery := "SELECT id, name, email FROM users"
	rows, _ := db.Query(getAllUsersQuery)

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func getUserHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userId, err := strconv.Atoi(idStr)

	user, err := GetUser(db, userId)
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUser(db *sql.DB, id int) (*User, error) {
	getUserQuery := "SELECT * FROM users WHERE id=?"
	row := db.QueryRow(getUserQuery, id)

	user := &User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func updateUserHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userId, err := strconv.Atoi(idStr)

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)

	UpdateUser(db, userId, user.Name, user.Email)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "User updated successfully")
}

func UpdateUser(db *sql.DB, id int, name string, email string) error {
	updateUserQuery := "UPDATE users SET name=?, email=? WHERE id=?"
	_, err := db.Exec(updateUserQuery, name, email, id)
	if err != nil {
		return err
	}
	return nil
}

func deleteUserHandle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	vars := mux.Vars(r)
	idStr := vars["id"]

	userId, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}

	user := DeleteUser(db, userId)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(db *sql.DB, id int) error {
	getUserQuery := "DELETE FROM users WHERE id=?"
	_, err := db.Exec(getUserQuery, id)
	if err != nil {
		return err
	}
	return nil
}
