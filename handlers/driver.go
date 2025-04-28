package handlers

import (
	"f1-statshub/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetAllDrivers maneja la solicitud GET para obtener todos los pilotos
func GetAllDrivers(c *gin.Context) {
	drivers, err := services.GetAllDrivers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener los pilotos"})
		return
	}
	var result []map[string]interface{}
	for _, driver := range drivers {
		result = append(result, map[string]interface{}{
			"first_name":  driver.FirstName,
			"last_name":   driver.LastName,
			"team_name":   driver.TeamName,
			"country_code": driver.CountryCode,
		})
	}

	c.JSON(http.StatusOK, result)
}

/// GetDriverDetail maneja GET /api/corredor/detalle/:id
func GetDriverDetail(c *gin.Context) {
	driverIDStr := c.Param("id")
	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de piloto inválido"})
		return
	}

	// Primera consulta: obtener performance_summary
	performanceSummary, err := services.GetDriverPerformanceSummary(driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el resumen de rendimiento"})
		return
	}

	// Segunda consulta: obtener race_results
	raceResults, err := services.GetDriverRaceResults(driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener resultados de las carreras"})
		return
	}

	// Armamos la respuesta final
	response := gin.H{
		"driver_id": driverID,
		"performance_summary": performanceSummary,
		"race_results": raceResults,
	}

	c.JSON(http.StatusOK, response)
}

