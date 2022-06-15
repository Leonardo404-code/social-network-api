package main

import (
	"api/src/config"
	"api/src/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	config.Load()

	r := routes.GenerateRoutes()

	log.Println("Api Running at :3001 port")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r))
}
