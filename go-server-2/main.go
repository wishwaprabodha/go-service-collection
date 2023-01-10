package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/wishwaprabodha/go-server/router"
	"github.com/wishwaprabodha/go-server/service"
	"log"
	"net/http"
)

func main() {
	log.Println("Server Starting")
	service.InitBooks()
	log.Println(http.ListenAndServe(":8002", router.StartRoutes()))
}
