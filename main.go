package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func TestEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Testing 1 2 3 ..."))
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type Star struct {
	ID          uint
	Name        string `gorm:"unique"`
	Description string
	URL         string
}

type App struct {

	DB *gorm.DB
	Router *mux.Router
}

func (a *App) Initialize(dbURI string) {

	//
	DB, err := gorm.Open(sqlite.Open(dbURI), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	a.DB = DB
	a.Router = mux.NewRouter()

	//Migrate the schema
	DB.AutoMigrate(&Product{}, &Star{})

}

func main() {
	//init db
	a := &App{}
	a.Initialize(":memory:")
	//a.Initialize("test.db")

	a.DB.Create(&Star{Name: "vb"})
	a.DB.Create(&Star{Name: "vb2"})




	r := a.Router

	r.HandleFunc("/", a.defaultRouter)

	r.HandleFunc("/bye", byeRouter)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)


	defer cleanDB(a.DB)

}

func byeRouter(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Bye")
	w.WriteHeader(http.StatusOK)
}

func (a *App)  defaultRouter(w http.ResponseWriter, r *http.Request) {
	var star Star
	a.DB.First(&star, "name = ?", "vb")
	fmt.Fprintf(w, "Hi")
	w.Write([]byte(star.Name))
}

//TODO: defer not the best approach here, . use shutdown hooks
func cleanDB(DB *gorm.DB) {
	db, _ := DB.DB()
	db.Close()
	println("Cleaned up DB")
	time.Sleep(time.Duration(time.Second*5))
}