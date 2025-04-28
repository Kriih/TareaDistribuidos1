package handlers

import (
	"f1-statshub/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetRaces maneja GET /api/carrera
func GetRaces(c *gin.Context) {
	races, err := services.GetRaces()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener carreras"})
		return
	}
	c.JSON(http.StatusOK, races)
}

// GetRaceDetail maneja GET /api/carrera/detalle/:id
func GetRaceDetail(c *gin.Context) {
	// Capturamos el id de la carrera desde la URL
	raceIDStr := c.Param("id")
	raceID, err := strconv.Atoi(raceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de carrera inválido"})
		return
	}

	// Llamamos al service que maneja la lógica
	raceDetail, err := services.GetRaceDetail(raceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener detalle de carrera"})
		return
	}

	c.JSON(http.StatusOK, raceDetail)
}
