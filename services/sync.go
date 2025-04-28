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

func SyncDrivers() error {
	for sessionKey, driverNumbers := range sessionDrivers {
		url := fmt.Sprintf("https://api.openf1.org/v1/drivers?session_key=%d", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error en request de drivers: %v", err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		var allDrivers []models.Driver
		if err := json.Unmarshal(body, &allDrivers); err != nil {
			return fmt.Errorf("Error parseando JSON de drivers: %v", err)
		}

		for _, d := range allDrivers {
			if contains(driverNumbers, d.DriverNumber) {
				insertDriver(d)
			}
		}
	}
	fmt.Println("Drivers sincronizados.")
	return nil
}

func SyncSessions() error {
	url := "https://api.openf1.org/v1/sessions?session_name=Race&year=2024"
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error al hacer request de sesiones: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var sessions []models.Session
	if err := json.Unmarshal(body, &sessions); err != nil {
		return fmt.Errorf("Error parseando JSON de sesiones: %v", err)
	}

	for _, s := range sessions {
		insertSession(s)
	}
	fmt.Println("Sesiones sincronizadas.")
	return nil
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

func SyncPositions() error {
	rows, err := db.DB.Query("SELECT session_key FROM sessions")
	if err != nil {
		return fmt.Errorf("Error al obtener sesiones: %v", err)
	}
	defer rows.Close()

	var keys []int
	for rows.Next() {
		var sessionKey int
		if err := rows.Scan(&sessionKey); err == nil {
			keys = append(keys, sessionKey)
		}
	}
	rows.Close()

	for _, sessionKey := range keys {
		url := fmt.Sprintf("https://api.openf1.org/v1/position?session_key=%d", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error request posiciones sesión %d: %v", sessionKey, err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var positions []models.Position
		if err := json.Unmarshal(body, &positions); err != nil {
			return fmt.Errorf("Error parseando posiciones sesión %d: %v", sessionKey, err)
		}

		for _, p := range positions {
			insertPosition(p)
		}
	}

	fmt.Println("Posiciones sincronizadas.")
	return nil
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

func SyncLaps() error {
	rows, err := db.DB.Query("SELECT session_key FROM sessions")
	if err != nil {
		return fmt.Errorf("Error al obtener sesiones: %v", err)
	}
	defer rows.Close()

	var sessionKeys []int
	for rows.Next() {
		var sessionKey int
		if err := rows.Scan(&sessionKey); err == nil {
			sessionKeys = append(sessionKeys, sessionKey)
		}
	}
	rows.Close()

	for _, sessionKey := range sessionKeys {
		url := fmt.Sprintf("https://api.openf1.org/v1/laps?session_key=%d", sessionKey)
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error request laps sesión %d: %v", sessionKey, err)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var laps []models.Lap
		if err := json.Unmarshal(body, &laps); err != nil {
			return fmt.Errorf("Error parseando laps sesión %d: %v", sessionKey, err)
		}

		for _, l := range laps {
			insertLap(l)
		}
	}
	fmt.Println("Vueltas sincronizadas.")
	return nil
}

func insertLap(l models.Lap) {
	stmt := `
		INSERT OR IGNORE INTO laps (
			driver_number, session_key, lap_number,
			lap_duration, duration_sector_1, duration_sector_2, duration_sector_3,
			st_speed, date_start
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.DB.Exec(stmt,
		l.DriverNumber, l.SessionKey, l.LapNumber,
		l.LapDuration, l.Sector1, l.Sector2, l.Sector3,
		l.StSpeed, l.DateStart,
	)
	if err != nil {
		fmt.Printf("Error insertando vuelta: %v\n", err)
	}
}
