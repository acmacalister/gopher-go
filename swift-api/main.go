package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/hex"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/codegangsta/negroni"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/mholt/binding"
	"gopkg.in/unrolled/render.v1"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type DBHandler struct {
	db *gorm.DB
	r  *render.Render
}

type Guitar struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Brand     string    `json:"brand"`
	Year      string    `json:"year"`
	Price     int64     `json:"price"`
	Color     string    `json:"color"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Our form values we need for updating/creating a guitar.
type GuitarForm struct {
	Name  string
	Brand string
	Year  string
	Price int64
	Color string
}

//Our User to auth our people
type User struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name"`
	PasswordDigest string    `json:"password_digest"`
	ImageUrl       string    `json:"image_url"`
	AuthToken      string    `json:"auth_token"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// to do some validation on our input fields. File is done separately.
func (gf *GuitarForm) FieldMap() binding.FieldMap {
	return binding.FieldMap{
		&gf.Name: binding.Field{
			Form:     "name",
			Required: true,
		},
		&gf.Brand: binding.Field{
			Form:     "brand",
			Required: true,
		},
		&gf.Year: binding.Field{
			Form:     "year",
			Required: true,
		},
		&gf.Price: binding.Field{
			Form:     "price",
			Required: true,
		},
		&gf.Color: binding.Field{
			Form:     "color",
			Required: true,
		},
	}
}

const (
	defaultPerPage = 30
)

func main() {
	// setup db. We would normally load this out of a config file,
	// but for this example, it is hardset. See gist at end of article for config example.
	db, err := gorm.Open("mysql", "root@/guitarstore?parseTime=true")

	if err != nil {
		log.Fatal(skittles.BoldRed(err))
	}
	db.LogMode(true) // This would be off in production.
	defer db.Close()
	db.AutoMigrate(&Guitar{}) // nice for development, but I would probably just write a SQL script to do this.
	db.AutoMigrate(&User{})

	r := render.New(render.Options{})
	h := DBHandler{db: &db, r: r}

	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/create", h.createUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", h.loginUserHandler).Methods("POST")

	// setup a basic CRUD/REST API for our guitar store.
	router := mux.NewRouter()
	router.HandleFunc("/guitars", h.guitarsIndexHandler).Methods("GET")
	router.HandleFunc("/guitars", h.guitarsCreateHandler).Methods("POST")
	router.HandleFunc("/guitars/{id:[0-9]+}", h.guitarsShowHandler).Methods("GET")
	router.HandleFunc("/guitars/{id:[0-9]+}", h.guitarsUpdateHandler).Methods("PUT", "PATCH")
	router.HandleFunc("/guitars/{id:[0-9]+}", h.guitarsDestroyHandler).Methods("DELETE")

	//auth the guitar routes
	authRouter.Handle("/guitars", negroni.New(
		negroni.HandlerFunc(h.authHandler),
		negroni.Wrap(router),
	))

	n := negroni.Classic()
	n.UseHandler(authRouter)
	n.Run(":8080")
}

// create a new user
func (h *DBHandler) createUserHandler(rw http.ResponseWriter, req *http.Request) {
	// Get the form values out of the POST request.
	name := req.FormValue("name")
	password := req.FormValue("password")
	imageUrl := req.FormValue("imageUrl")

	// Generate a hashed password from bcrypt.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	count := 16
	b := make([]byte, count)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < count; i++ {
		b[i] = byte(rand.Intn(count))
	}
	token := hex.EncodeToString(b)
	// Stick that in our users table of our db.
	user := User{Name: name, PasswordDigest: string(hashedPass), ImageUrl: imageUrl, AuthToken: token, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	h.db.Save(&user)
	user.PasswordDigest = "" //we don't need to expose that to the user
	h.r.JSON(rw, http.StatusOK, &user)
}

//allows an existing user to login
func (h *DBHandler) loginUserHandler(rw http.ResponseWriter, req *http.Request) {
	// Get the form values out of the POST request.
	name := req.FormValue("name")
	password := req.FormValue("password")

	user := User{}
	h.db.Where("name = ?", name).First(&user) //in production code, we would of course validate this before running a where statement on a raw value
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password)); err != nil {
		log.Println("login error: ", err)
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}
	user.PasswordDigest = "" //we don't need to expose that to the user
	h.r.JSON(rw, http.StatusOK, &user)
}

//middleware that checks to make sure the authToken is a valid user
func (h *DBHandler) authHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token := r.FormValue("auth_token")
	user := User{}
	h.db.Where("auth_token = ?", token).First(&user) //in production code, we would of course validate this before running a where statement on a raw value
	if user.Name == "" {
		http.Error(rw, "Not authorized", http.StatusUnauthorized)
		return
	}
	next(rw, r)
}

// our guitar routes.

// guitarsIndexHandler returns all our guitars out of the db in a paginated fashion.
func (h *DBHandler) guitarsIndexHandler(rw http.ResponseWriter, req *http.Request) {
	page := getPage(req) - 1
	perPage := getPerPage(req)
	offset := perPage * page
	var guitars []Guitar
	h.db.Limit(perPage).Offset(offset).Find(&guitars)
	if guitars == nil {
		h.r.JSON(rw, http.StatusOK, "[]") // If we have no guitars, just return an empty array, instead of null.
	} else {
		h.r.JSON(rw, http.StatusOK, &guitars)
	}
}

// guitarsShowHandler returns a single guitar from the db.
func (h *DBHandler) guitarsShowHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	guitar := Guitar{}
	h.db.First(&guitar, id)
	h.r.JSON(rw, http.StatusOK, &guitar)
}

// guitarsCreateHandler inserts a new guitar into the db.
func (h *DBHandler) guitarsCreateHandler(rw http.ResponseWriter, req *http.Request) {
	h.guitarsEdit(rw, req, 0)
}

// guitarsUpdateHandler updates a guitar in the db.
func (h *DBHandler) guitarsUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	h.guitarsEdit(rw, req, id)
}

// guitarsEdit is shared between the create and update handler, since they share most of the logic.
func (h *DBHandler) guitarsEdit(rw http.ResponseWriter, req *http.Request, id int64) {
	guitarForm := GuitarForm{}
	if err := binding.Bind(req, &guitarForm); err.Handle(rw) {
		return
	}

	// normally we would upload to S3, but for this demo, we will just write to disk. See this gist for S3 upload code.
	upload, header, err := req.FormFile("file")
	if err != nil {
		h.r.JSON(rw, http.StatusBadRequest, map[string]string{"error": "bad file upload."})
		return
	}
	file, err := os.Create(fmt.Sprintf("public/%s", header.Filename)) // we would normally need to generate unique filenames.
	if err != nil {
		h.r.JSON(rw, http.StatusInternalServerError, map[string]string{"error": "system error occured"})
		return
	}
	io.Copy(file, upload) // write the uploaded file to disk.
	imageUrl := fmt.Sprintf("http://localhost:8080/%s", header.Filename)

	guitar := Guitar{Id: id, Name: guitarForm.Name, Brand: guitarForm.Brand, Year: guitarForm.Year,
		Price: guitarForm.Price, Color: guitarForm.Color, ImageUrl: imageUrl, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	h.db.Save(&guitar)
	h.r.JSON(rw, http.StatusOK, &guitar)
}

// guitarsDestroyHandler deletes a guitar from the db.
func (h *DBHandler) guitarsDestroyHandler(rw http.ResponseWriter, req *http.Request) {
	id := getId(req)
	guitar := Guitar{}
	h.db.Delete(&guitar, id)
	h.r.JSON(rw, http.StatusOK, &guitar)
}

// getId parses our id out of the url.
func getId(req *http.Request) int64 {
	vars := mux.Vars(req)
	idString := vars["id"]
	id, err := strconv.ParseInt(idString, 10, 0)
	if err != nil {
		log.Println(skittles.BoldRed(err))
	}
	return id
}

// getPage returns the page param from the url query.
func getPage(req *http.Request) int {
	return parseQueryValues(req, "page")
}

// getPerPage returns the per_page param from the url query.
func getPerPage(req *http.Request) int {
	perPage := parseQueryValues(req, "per_page")
	if perPage == 0 {
		return defaultPerPage
	}
	return perPage
}

// parseQueryValues shared parser for page & per_page url query.
func parseQueryValues(req *http.Request, value string) int {
	vals := req.URL.Query()
	val := vals[value]
	if val != nil {
		v, err := strconv.ParseInt(val[0], 10, 0)
		if err != nil {
			log.Println(skittles.BoldRed(err))
			return 0
		}
		return int(v)
	}
	return 0
}
