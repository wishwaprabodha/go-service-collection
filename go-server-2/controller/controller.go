package controller

import (
	"encoding/json"
	"fmt"
	"github.com/dolthub/vitess/go/vt/log"
	"github.com/gorilla/mux"
	"github.com/wishwaprabodha/go-service-collection/go-server-2/model"
	"github.com/wishwaprabodha/go-service-collection/go-server-2/service"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Server Started & Running...")
	json.NewEncoder(w).Encode("all is well")
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book service.Book
	books := book.GetBooks()
	json.NewEncoder(w).Encode(books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var book service.Book
	searchBook := book.GetBook(params["id"])
	json.NewEncoder(w).Encode(searchBook)
}

func AddBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book service.Book
	error := json.NewDecoder(r.Body).Decode(&book)
	if error != nil {
		http.Error(w, "Decode Error", 400)
	}
	book.AddBook(book)
	json.NewEncoder(w).Encode(book)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book service.Book
	params := mux.Vars(r)
	error := json.NewDecoder(r.Body).Decode(&book)
	if error != nil {
		http.Error(w, "Decode Error", 400)
	}
	book.UpdateBook(params["id"], book)
	json.NewEncoder(w).Encode(book)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book service.Book
	params := mux.Vars(r)
	book.DeleteBook(params["id"])
	json.NewEncoder(w).Encode(book)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Decode Error", 400)
	}
	key, detaErr := service.CreateUserDeta(&user)
	if detaErr != nil {
		log.Fatal("Error: ", detaErr)
	}
	log.Info("key: ", key)
	err, _ = service.CreateUser(&user)
	if err != nil {
		return
	}
	errJson := json.NewEncoder(w).Encode(user)
	if errJson != nil {
		log.Fatal("Error: ", errJson)
	}
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	err, user := service.GetUserByEmail(params["email"])
	fmt.Println(user)
	if err != nil {
		return
	}
	errJson := json.NewEncoder(w).Encode(&user)
	if errJson != nil {
		log.Fatal("Error: ", errJson)
	}
}
