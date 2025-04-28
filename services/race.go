package services

import (
	"f1-statshub/db"
)

// GetRaces obtiene todas las carreras (sessions tipo "Race" año 2024)
func GetRaces() ([]map[string]interface{}, error) {
	rows, err := db.DB.Query(`
		SELECT session_key, country_name, date_start, year, circuit_short_name
		FROM sessions
		WHERE session_type = 'Race' AND year = 2024
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var races []map[string]interface{}

	for rows.Next() {
		var sessionKey int
		var countryName, dateStart, circuitShortName string
		var year int

		err := rows.Scan(&sessionKey, &countryName, &dateStart, &year, &circuitShortName)
		if err == nil {
			races = append(races, map[string]interface{}{
				"session_key":        sessionKey,
				"country_name":       countryName,
				"date_start":         dateStart,
				"year":               year,
				"circuit_short_name": circuitShortName,
			})
		}
	}

	return races, nil
}

// GetRaceDetail obtiene todo el resumen de una carrera
func GetRaceDetail(raceID int) (map[string]interface{}, error) {
	var countryName, dateStart, circuitShortName string
	var year int

	err := db.DB.QueryRow(`
		SELECT country_name, date_start, circuit_short_name, year
		FROM sessions
		WHERE session_key = ?
	`, raceID).Scan(&countryName, &dateStart, &circuitShortName, &year)
	if err != nil {
		return nil, err
	}

	// Resultados de la carrera
	type PodiumResult struct {
		Position interface{} `json:"position"`
		Driver   string      `json:"driver"`
		Team     string      `json:"team"`
		Country  string      `json:"country"`
	}

	rows, err := db.DB.Query(`
		SELECT p.position, d.first_name || ' ' || d.last_name, d.team_name, d.country_code
		FROM positions p
		JOIN drivers d ON p.driver_number = d.driver_number
		WHERE p.session_key = ?
		ORDER BY p.position ASC
	`, raceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []PodiumResult

	for rows.Next() {
		var r PodiumResult
		err := rows.Scan(&r.Position, &r.Driver, &r.Team, &r.Country)
		if err == nil {
			results = append(results, r)
		}
	}

	var podium []PodiumResult
	var last PodiumResult

	if len(results) >= 3 {
		podium = results[:3]
		last = results[len(results)-1]
	}

	// Fastest Lap
	var fastDriver string
	var totalTime, sector1, sector2, sector3 float64

	err = db.DB.QueryRow(`
		SELECT d.first_name || ' ' || d.last_name, MIN(l.lap_duration),
		       MIN(l.duration_sector_1), MIN(l.duration_sector_2), MIN(l.duration_sector_3)
		FROM laps l
		JOIN drivers d ON l.driver_number = d.driver_number
		WHERE l.session_key = ?
	`, raceID).Scan(&fastDriver, &totalTime, &sector1, &sector2, &sector3)

	if err != nil {
		fastDriver = ""
	}

	// Max Speed
	var maxSpeedDriver string
	var maxSpeed float64

	err = db.DB.QueryRow(`
		SELECT d.first_name || ' ' || d.last_name, MAX(l.st_speed)
		FROM laps l
		JOIN drivers d ON l.driver_number = d.driver_number
		WHERE l.session_key = ?
	`, raceID).Scan(&maxSpeedDriver, &maxSpeed)

	if err != nil {
		maxSpeedDriver = ""
	}

	// Construir la respuesta
	return map[string]interface{}{
		"race_id":            raceID,
		"country_name":       countryName,
		"date_start":         dateStart,
		"year":               year,
		"circuit_short_name": circuitShortName,
		"results": []map[string]interface{}{
			{"position": podium[0].Position, "driver": podium[0].Driver, "team": podium[0].Team, "country": podium[0].Country},
			{"position": podium[1].Position, "driver": podium[1].Driver, "team": podium[1].Team, "country": podium[1].Country},
			{"position": podium[2].Position, "driver": podium[2].Driver, "team": podium[2].Team, "country": podium[2].Country},
			{"position": "Ultimo", "driver": last.Driver, "team": last.Team, "country": last.Country},
		},
		"fastest_lap": map[string]interface{}{
			"driver":     fastDriver,
			"total_time": totalTime,
			"sector_1":   sector1,
			"sector_2":   sector2,
			"sector_3":   sector3,
		},
		"max_speed": map[string]interface{}{
			"driver":    maxSpeedDriver,
			"speed_kmh": maxSpeed,
		},
	}, nil
}
