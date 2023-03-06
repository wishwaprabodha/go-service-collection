package router

import (
	"github.com/gorilla/mux"
	"github.com/wishwaprabodha/go-service-collection/go-server-1/controller"
)

func StartRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", controller.Index).Methods("GET")
	router.HandleFunc("/api/books", controller.GetBooks).Methods("GET")
	router.HandleFunc("/api/book/{id}", controller.GetBook).Methods("GET")
	router.HandleFunc("/api/book", controller.AddBook).Methods("POST")
	router.HandleFunc("/api/user", controller.AddUser).Methods("POST")
	router.HandleFunc("/api/user/{email}", controller.GetUserByEmail).Methods("GET")
	router.HandleFunc("/api/book/{id}", controller.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/book/{id}", controller.DeleteBook).Methods("DELETE")
	return router
}
