package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/wishwaprabodha/go-service-collection/go-server-2/router"
	"github.com/wishwaprabodha/go-service-collection/go-server-2/service"
	"log"
	"net/http"
)

func main() {
	log.Println("Server Starting")
	service.InitBooks()
	log.Println(http.ListenAndServe(":8002", router.StartRoutes()))
}
