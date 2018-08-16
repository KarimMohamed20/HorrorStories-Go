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

type User struct {
	gorm.Model
	Name  string
	Story string
	Shortcut string
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func newUser(w http.ResponseWriter, r *http.Request) {
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
	db.Create(&User{Name: name, Story: story,Shortcut:shortcut})
	fmt.Fprintf(w, "New User Successfully Created")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]

	var user User
	db.Where("name = ?", name).Find(&user)
	db.Delete(&user)

	fmt.Fprintf(w, "Successfully Deleted User")
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	vars := mux.Vars(r)
	name := vars["name"]
	story := vars["story"]

	var user User
	db.Where("name = ?", name).Find(&user)

	user.Story = story

	db.Save(&user)
	fmt.Fprintf(w, "Successfully Updated User")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/stories", allUsers).Methods("GET")
	myRouter.HandleFunc("/story/{name}", deleteUser).Methods("DELETE")
	myRouter.HandleFunc("/story/{name}/{story}/{shortcut}", updateUser).Methods("PUT")
	myRouter.HandleFunc("/story/{name}/{story}/{shortcut}", newUser).Methods("POST")
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
	db.AutoMigrate(&User{})
}

func main() {
	fmt.Println("Horror Stories")

	initialMigration()
	// Handle Subsequent requests
	handleRequests()
}
