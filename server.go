package main

import (
	"f1-statshub/db"
	"f1-statshub/services"
	"f1-statshub/routes"
	"fmt"
	"log"
	"net/http"
)

func StartServer() {
	r := routes.NewRouter()

	// Iniciar la api
	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


func main() {
	db.InitDB()
	services.SyncDrivers()
	services.SyncSessions()
	services.SyncPositions()
	services.SyncLaps()

	StartServer()
}
