package main

import (
	"github.com/hoisie/mustache"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Guitar struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Year      int64     `json:"year"`
	Price     int64     `json:"price"`
	Color     string    `json:"color"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	g := Guitar{Id: 1, Name: "Les Paul", Brand: "Gibson", Year: 1966, Price: 3500, Color: "Sunburst Cherry",
		CreatedAt: time.Now(), UpdatedAt: time.Now()} // normally we would pull this out of a db. See the API article on how to.
	http.HandleFunc("/", g.templateHandler)
	http.HandleFunc("/mustache", g.mustacheHandler)
	http.ListenAndServe(":8080", nil)
}

func (g *Guitar) templateHandler(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("guitar.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(rw, g) // notice guitar is a pointer.
}

func (g *Guitar) mustacheHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(mustache.RenderFile("guitar.mustache", g))) // notice guitar is a pointer.
}
