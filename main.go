package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Story struct {
	gorm.Model
	Name  string
	Story string
	Shortcut string
}

func allStories(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var stories []Story
	db.Find(&stories)
	json.NewEncoder(w).Encode(stories)
}

func newStories(w http.ResponseWriter, r *http.Request) {
	fmt.Println("New User Endpoint Hit")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	story := vars["story"]
	shortcut := vars["shortcut"]

	fmt.Println(name)
	fmt.Println(story)
	fmt.Println(shortcut)
	db.Create(&Story{Name: name, Story: story,Shortcut:shortcut})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteStories(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var stories Story
	db.Where("name = ?", name).Find(&stories)
	db.Delete(&stories)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateStories(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	story := vars["story"]

	var stories Story
	db.Where("name = ?", name).Find(&stories)

	stories.Story = story

	db.Save(&stories)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/stories", allStories).Methods("GET")
	myRouter.HandleFunc("/story/{name}", deleteStories).Methods("DELETE")
	myRouter.HandleFunc("/story/{name}/{story}/{shortcut}", updateStories).Methods("PUT")
	myRouter.HandleFunc("/story/{name}/{story}/{shortcut}", newStories).Methods("POST")
	log.Fatal(http.ListenAndServe(":8010", myRouter))
}

func initialMigration() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Story{})
}

func main() {
	fmt.Println("Horror Stories")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
