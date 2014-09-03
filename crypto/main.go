package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

// Contains our sql connection.
type DBHandler struct {
	db *sql.DB
}

func main() {
	// Open SQL connection to db.
	db, err := sql.Open("mysql", "root@/tester")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Pass it to the handler.
	h := DBHandler{db: db}

	// A couple HTTP routes.
	http.HandleFunc("/create", h.CreateHandler)
	http.HandleFunc("/auth", h.AuthHandler)
	http.ListenAndServe(":8081", nil)
}

func (h *DBHandler) CreateHandler(rw http.ResponseWriter, req *http.Request) {
	// Get the form values out of the POST request.
	name := req.FormValue("name")
	password := req.FormValue("password")

	// Generate a hashed password from bcrypt.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	// Stick that in our users table of our db.
	_, err = h.db.Query("INSERT INTO users (name, password_digest, created_at, updated_at) VALUES(?,?,?,?)",
		name, hashedPass, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
	}

	// Write a silly message back to the client.
	rw.Write([]byte("Created user!"))
}

func (h *DBHandler) AuthHandler(rw http.ResponseWriter, req *http.Request) {
	// Get the form values out of the POST request.
	name := req.FormValue("name")
	password := req.FormValue("password")

	// Find the user by his name and get the password_digest we generated in the create method out.
	var digest string
	if err := h.db.QueryRow("SELECT password_digest FROM users WHERE name = ?", name).Scan(&digest); err != nil {
		log.Fatal(err)
	}

	// Compare that password_digest to our password we got from the form value.
	// If the error is not equal to nil, we know the auth failed. If there is no error, it
	// was successful.
	if err := bcrypt.CompareHashAndPassword([]byte(digest), []byte(password)); err != nil {
		rw.Write([]byte("auth failure..."))
	} else {
		rw.Write([]byte("auth successful!"))
	}
}
