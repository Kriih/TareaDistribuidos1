package main

import (
	"f1-statshub/db"
	"f1-statshub/routes"
	"f1-statshub/services"
	"fmt"
)

func StartServer() {

	// Obtener router configurado (GIN)
	r := routes.NewRouter()

	// Iniciar servidor con GIN directamente
	fmt.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}

func main() {
	db.InitDB()
	services.SyncDrivers()
	services.SyncSessions()
	services.SyncPositions()
	services.SyncLaps()

	StartServer()
}
