package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type DBHandler struct {
	db *sql.DB
}

func main() {
	db, err := sql.Open("mysql", "root@/mobile_lsfilter_dev")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	h := DBHandler{db: db}

	http.HandleFunc("/", h.HomeHandler)
	http.ListenAndServe(":8081", nil)
}

func (h *DBHandler) HomeHandler(rw http.ResponseWriter, req *http.Request) {
	name := "dalton"
	password := []byte("test")
	var digest string
	if err := h.db.QueryRow("SELECT password_digest FROM users WHERE name = ?", name).Scan(&digest); err != nil {
		log.Fatal(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(digest), password); err != nil {
		rw.Write([]byte("auth failure..."))
	} else {
		rw.Write([]byte("auth successful!"))
	}
}
