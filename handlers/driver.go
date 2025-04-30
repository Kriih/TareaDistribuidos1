package handlers

import (
	"f1-statshub/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
			"driver_number": driver.DriverNumber,
			"team_name":   driver.TeamName,
			"country_code": driver.CountryCode,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetDriverDetail(c *gin.Context) {
	driverIDStr := c.Param("id")
	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de piloto inválido"})
		return
	}

	performanceSummary, err := services.GetDriverPerformanceSummary(driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el resumen de rendimiento"})
		return
	}

	raceResults, err := services.GetDriverRaceResults(driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener resultados de las carreras"})
		return
	}

	response := gin.H{
		"driver_id": driverID,
		"performance_summary": performanceSummary,
		"race_results": raceResults,
	}

	c.JSON(http.StatusOK, response)
}

