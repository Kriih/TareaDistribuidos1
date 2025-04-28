package main

import (
	"f1-statshub/db"
	"f1-statshub/routes"
	"f1-statshub/services"
	"fmt"
	"log"
	"net/http"
	"time"
)

func StartServer() {
	r := routes.NewRouter()

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	db.InitDB()

	for {
		// Ejecutar las sincronizaciones
		if err := services.SyncDrivers(); err != nil {
			log.Printf("Error sincronizando drivers: %v. Reintentando en 3 segundos...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if err := services.SyncSessions(); err != nil {
			log.Printf("Error sincronizando sessions: %v. Reintentando en 3 segundos...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if err := services.SyncPositions(); err != nil {
			log.Printf("Error sincronizando positions: %v. Reintentando en 3 segundos...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		if err := services.SyncLaps(); err != nil {
			log.Printf("Error sincronizando laps: %v. Reintentando en 3 segundos...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	StartServer()
}
