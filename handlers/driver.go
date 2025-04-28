package handlers

import (
	"f1-statshub/services"
	"github.com/gin-gonic/gin"
	"net/http"
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
