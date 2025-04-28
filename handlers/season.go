package handlers

import (
	"f1-statshub/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSeasonSummary maneja GET /api/temporada/resumen
func GetSeasonSummary(c *gin.Context) {
	summary, err := services.GetSeasonSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener resumen de temporada"})
		return
	}
	c.JSON(http.StatusOK, summary)
}
