package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Participant struct {
	gorm.Model
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique"`
}

type Issue struct {
	gorm.Model
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique"`
}

type ParticipantIssueRating struct {
	gorm.Model
	ParticipantID uint `gorm:"primaryKey;autoIncrement:false"`
	Participant   Participant
	IssueID       uint `gorm:"primaryKey;autoIncrement:false"`
	Issue         Issue
	Value         float32
}

type App struct {
	DB     *gorm.DB
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
	DB.Debug().AutoMigrate(&Participant{}, &Issue{}, &ParticipantIssueRating{})

}

func main() {
	//init db
	a := &App{}
	a.Initialize(":memory:")
	//a.Initialize("test.db")

	participant := Participant{Name: "vinay"}
	a.DB.Debug().Create(&participant)
	issue := Issue{Name: "ducks in the pool"}
	a.DB.Debug().Create(&issue)
	issue2 := Issue{Name: "pigs in the garage"}
	a.DB.Debug().Create(&issue2)
	a.DB.Debug().Create(&ParticipantIssueRating{
		Participant: participant,
		Issue:       issue,
		Value:       1.4,
	})
	a.DB.Debug().Create(&ParticipantIssueRating{
		Participant: participant,
		Issue:       issue2,
		Value:       5,
	})
	//a.DB.Debug().Create(&ParticipantIssueRating{
	//	Participant: participant,
	//	Issue:       issue2,
	//	Value:       5,
	//})

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

func (a *App) defaultRouter(w http.ResponseWriter, r *http.Request) {
	var participant Participant
	a.DB.First(&participant, "name = ?", "vinay")
	fmt.Fprintf(w, "Hi")
	//w.Write([]byte(star.Name))
}

//TODO: defer not the best approach here, . use shutdown hooks
func cleanDB(DB *gorm.DB) {
	db, _ := DB.DB()
	db.Close()
	println("Cleaned up DB")
	time.Sleep(time.Duration(time.Second * 5))
}
