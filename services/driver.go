package services

import (
	"f1-statshub/db"
	"f1-statshub/models"
	"log"
)

// GetAllDrivers obtiene todos los pilotos desde la base de datos
func GetAllDrivers() ([]models.Driver, error) {
	var drivers []models.Driver

	// Consulta para obtener todos los drivers
	rows, err := db.DB.Query("SELECT driver_number, first_name, last_name, name_acronym, team_name, country_code FROM drivers")
	if err != nil {
		log.Printf("Error ejecutando consulta: %v", err)
		return nil, err
	}
	defer rows.Close()

	// Iteramos sobre los resultados
	for rows.Next() {
		var driver models.Driver
		if err := rows.Scan(&driver.DriverNumber, &driver.FirstName, &driver.LastName, &driver.NameAcronym, &driver.TeamName, &driver.CountryCode); err != nil {
			log.Printf("Error al escanear fila: %v", err)
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	// Verificamos si hubo errores en la iteración
	if err := rows.Err(); err != nil {
		log.Printf("Error al iterar sobre las filas: %v", err)
		return nil, err
	}

	return drivers, nil
}

func GetDriverRaceResults(driverID int) ([]models.DriverDetail, error) {
	rows, err := db.DB.Query(`
		SELECT 
			positions.session_key,
			sessions.circuit_short_name,
			"GP de " || sessions.circuit_short_name AS race,
			positions.position,
			CASE WHEN MIN(laps.lap_duration) = laps.lap_duration THEN 1 ELSE 0 END AS fastest_lap,
			MAX(laps.st_speed) AS max_speed,
			MIN(laps.lap_duration) AS best_lap_duration
		FROM positions
		JOIN sessions ON positions.session_key = sessions.session_key
		LEFT JOIN laps ON positions.driver_number = laps.driver_number AND positions.session_key = laps.session_key
		WHERE positions.driver_number = ?
		GROUP BY positions.session_key
	`, driverID)

	if err != nil {
		log.Printf("Error ejecutando consulta de resultados: %v", err)
		return nil, err
	}
	defer rows.Close()

	var raceResults []models.DriverDetail
	for rows.Next() {
		var result models.DriverDetail
		var fastestLapInt int

		if err := rows.Scan(
			&result.SessionKey,
			&result.CircuitShortName,
			&result.Race,
			&result.Position,
			&fastestLapInt,
			&result.MaxSpeed,
			&result.BestLapDuration,
		); err != nil {
			log.Printf("Error escaneando fila de resultados: %v", err)
			return nil, err
		}

		result.FastestLap = fastestLapInt == 1
		raceResults = append(raceResults, result)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error al iterar resultados de carreras: %v", err)
		return nil, err
	}

	return raceResults, nil
}


func GetDriverPerformanceSummary(driverID int) (models.PerformanceSummary, error) {
	var summary models.PerformanceSummary

	query := `
		SELECT 
			COUNT(CASE WHEN position = 1 THEN 1 END) AS wins,
			COUNT(CASE WHEN position <= 3 THEN 1 END) AS top_3_finishes,
			MAX(laps.st_speed) AS max_speed
		FROM positions
		LEFT JOIN laps ON positions.driver_number = laps.driver_number AND positions.session_key = laps.session_key
		WHERE positions.driver_number = ?;
	`

	err := db.DB.QueryRow(query, driverID).Scan(&summary.Wins, &summary.Top3Finishes, &summary.MaxSpeed)
	if err != nil {
		log.Printf("Error ejecutando consulta de resumen: %v", err)
		return summary, err
	}

	return summary, nil
}
