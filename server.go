package main

import (
	"f1-statshub/db"
	"f1-statshub/services"
)

func main() {
	db.InitDB()
	services.SyncDrivers()
	services.SyncSessions()
	services.SyncPositions()
}
