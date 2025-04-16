package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"f1-statshub/db"
	"f1-statshub/models"
)

// Numeros requeridos según enunciado
var sessionDrivers = map[int][]int{
	9574: {1, 2, 3, 4, 10, 11, 14, 16, 18, 20, 22, 23, 24, 27, 31, 44, 55, 63, 77, 81},
	9636: {30, 50, 43},
}

func SyncDrivers() {
	for sessionKey, driverNumbers := range sessionDrivers {
		url := fmt.Sprintf("https://api.openf1.org/v1/drivers?session_key=%d", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error en request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var allDrivers []models.Driver
		if err := json.Unmarshal(body, &allDrivers); err != nil {
			fmt.Printf("Error parseando JSON: %v\n", err)
			continue
		}

		for _, d := range allDrivers {
			if contains(driverNumbers, d.DriverNumber) {
				insertDriver(d)
			}
		}
	}
	fmt.Println("Drivers sincronizados.")
}

func SyncSessions() {
	url := "https://api.openf1.org/v1/sessions?session_name=Race&year=2024"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error al hacer request de sesiones: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var sessions []models.Session
	if err := json.Unmarshal(body, &sessions); err != nil {
		fmt.Printf("Error parseando JSON de sesiones: %v\n", err)
		return
	}

	for _, s := range sessions {
		insertSession(s)
	}
	fmt.Println("Sesiones sincronizadas.")
}

func insertSession(s models.Session) {
	stmt := `
		INSERT OR IGNORE INTO sessions (
			session_key, session_name, session_type, location, country_name, year, circuit_short_name, date_start
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(stmt,
		s.SessionKey, s.SessionName, s.SessionType, s.Location,
		s.CountryName, s.Year, s.CircuitShortName, s.DateStart,
	)
	if err != nil {
		fmt.Printf("Error insertando sesión %d: %v\n", s.SessionKey, err)
	}
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func SyncPositions() {
	rows, err := db.DB.Query("SELECT session_key FROM sessions")
	if err != nil {
		fmt.Printf("Error al obtener sesiones: %v\n", err)
		return
	}
	defer rows.Close()

	var sessionKey int
	for rows.Next() {
		if err := rows.Scan(&sessionKey); err != nil {
			continue
		}

		url := fmt.Sprintf("https://api.openf1.org/v1/position?session_key=%d", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error request posiciones sesión %d: %v\n", sessionKey, err)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var positions []models.Position
		if err := json.Unmarshal(body, &positions); err != nil {
			fmt.Printf("Error parseando posiciones sesión %d: %v\n", sessionKey, err)
			continue
		}

		for _, p := range positions {
			stmt := `
				INSERT INTO positions (
					driver_number, session_key, position, date
				) VALUES (?, ?, ?, ?)
			`
			_, err := db.DB.Exec(stmt, p.DriverNumber, p.SessionKey, p.Position, p.Date)
			if err != nil {
				fmt.Printf("Error insertando posición sesión %d: %v\n", sessionKey, err)
			}
		}
	}

	fmt.Println("Posiciones sincronizadas.")
}

func insertPosition(p models.Position) {
	stmt := `
		INSERT INTO positions (
			driver_number, session_key, position, date
		) VALUES (?, ?, ?, ?)
	`
	_, err := db.DB.Exec(stmt, p.DriverNumber, p.SessionKey, p.Position, p.Date)
	if err != nil {
		fmt.Printf("Error insertando posición: %v\n", err)
	}
}

func insertDriver(d models.Driver) {
	stmt := `
		INSERT OR IGNORE INTO drivers (
			driver_number, first_name, last_name, name_acronym, team_name, country_code
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(stmt, d.DriverNumber, d.FirstName, d.LastName, d.NameAcronym, d.TeamName, d.CountryCode)
	if err != nil {
		fmt.Printf("Error insertando driver %d: %v\n", d.DriverNumber, err)
	}
}
