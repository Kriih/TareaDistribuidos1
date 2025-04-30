package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "proxy.db")
	if err != nil {
		log.Fatalf("Error abriendo base de datos: %v", err)
	}

	fmt.Println("Instalando/actualizando la base de datos...")
	clearTables()
	createTables()
	
}

func clearTables() {
	_, err := DB.Exec("DROP TABLE IF EXISTS drivers")
	if err != nil {
		log.Printf("Error eliminando tabla drivers: %v", err)
	}

	_, err = DB.Exec("DROP TABLE IF EXISTS sessions")
	if err != nil {
		log.Printf("Error eliminando tabla sessions: %v", err)
	}

	_, err = DB.Exec("DROP TABLE IF EXISTS positions")
	if err != nil {
		log.Printf("Error eliminando tabla positions: %v", err)
	}

	_, err = DB.Exec("DROP TABLE IF EXISTS laps")
	if err != nil {
		log.Printf("Error eliminando tabla laps: %v", err)
	}

	fmt.Println("Tablas eliminadas correctamente.")
}

func createTables() {
	createDriverTable := `
	CREATE TABLE IF NOT EXISTS drivers (
		driver_number INTEGER PRIMARY KEY,
		first_name TEXT,
		last_name TEXT,
		name_acronym TEXT,
		team_name TEXT,
		country_code TEXT
	);`

	createSessionTable := `
	CREATE TABLE IF NOT EXISTS sessions (
		session_key INTEGER PRIMARY KEY,
		session_name TEXT,
		session_type TEXT,
		location TEXT,
		country_name TEXT,
		year INTEGER,
		circuit_short_name TEXT,
		date_start TEXT
	);`

	createPositionTable := `
	CREATE TABLE IF NOT EXISTS positions (
		driver_number INTEGER,
		session_key INTEGER,
		position INTEGER,
		date TEXT
	);`

	createLapTable := `
	CREATE TABLE IF NOT EXISTS laps (
    driver_number INTEGER,
    session_key INTEGER,
    lap_number INTEGER,
    lap_duration REAL,
    duration_sector_1 REAL NULL,  
    duration_sector_2 REAL NULL,  
    duration_sector_3 REAL NULL,  
    st_speed REAL NULL,           
    date_start TEXT NULL          
	);`

	statements := []string{
		createDriverTable, createSessionTable, createPositionTable, createLapTable,
	}

	for _, stmt := range statements {
		_, err := DB.Exec(stmt)
		if err != nil {
			log.Fatalf("Error creando tabla: %v", err)
		}
	}

	fmt.Println("Tablas creadas correctamente.")
}
