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