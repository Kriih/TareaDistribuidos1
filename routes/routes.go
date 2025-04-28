package routes

import (
	"f1-statshub/handlers"

	"github.com/gin-gonic/gin"
)

// NewRouter crea y devuelve un enrutador configurado
func NewRouter() *gin.Engine {
	r := gin.Default()

	// Rutas para los usuarios
	r.GET("/api/carrera", handlers.GetRaces)
	r.GET("/api/carrera/detalle/:id", handlers.GetRaceDetail)
	r.GET("/api/temporada/resumen", handlers.GetSeasonSummary)

	// r.GET("/api/carrera/detalle/:id", handlers.GetRaceDetail)
	r.GET("/api/corredor", handlers.GetAllDrivers)
	return r
}
