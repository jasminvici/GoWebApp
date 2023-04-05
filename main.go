package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func main() {

	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/webappdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getAllUsers(db, w, r)
		case "POST":
			createUser(db, w, r)
		default:
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
		if err != nil {
			http.Error(w, "Invalid user ID.", 400)
			return
		}

		switch r.Method {
		case "GET":
			getUser(db, w, r, id)
		case "PUT":
			updateUser(db, w, r, id)
		case "DELETE":
			deleteUser(db, w, r, id)
		default:
			http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getAllUsers(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	sizeStr := r.URL.Query().Get("size")

	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		pageInt = 1
	}

	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil {
		sizeInt = 10
	}

	offset := (pageInt - 1) * sizeInt

	rows, err := db.Query("SELECT * FROM users LIMIT ?, ?", offset, sizeInt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	res, err := db.Exec("INSERT INTO users (name, email, age)	VALUES (?, ?, ?)", user.Name, user.Email, user.Age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	user.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getUser(db *sql.DB, w http.ResponseWriter, r *http.Request, id int) {
	row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found.", 404)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(db *sql.DB, w http.ResponseWriter, r *http.Request, id int) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	exists, err := userExists(db, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !exists {
		http.Error(w, "User not found.", 404)
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ?, age = ? WHERE id = ?", user.Name, user.Email, user.Age, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func userExists(db *sql.DB, id int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func deleteUser(db *sql.DB, w http.ResponseWriter, r *http.Request, id int) {

	exists, err := userExists(db, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if !exists {
		http.Error(w, "User not found.", 404)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, "User with ID %d deleted.", id)
}
