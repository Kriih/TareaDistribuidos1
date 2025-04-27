package routes

import (
	"f1-statshub/handlers"
	"github.com/gorilla/mux"
	
)

// NewRouter crea y devuelve un enrutador configurado
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Rutas para los usuarios
	r.HandleFunc("/api", handlers.HomeHandler).Methods("GET")
	return r
}
