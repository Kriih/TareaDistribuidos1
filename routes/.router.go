package routes

import (
	"f1-statshub/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

// NewRouter crea y devuelve un enrutador configurado
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Rutas para los usuarios
	r.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id}", handlers.GetUser).Methods("GET")
	r.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	// Rutas para los drivers
	r.HandleFunc("/drivers", handlers.GetAllDrivers).Methods("GET")
	r.HandleFunc("/drivers/{id}", handlers.GetDriver).Methods("GET")
	r.HandleFunc("/drivers", handlers.CreateDriver).Methods("POST")
	r.HandleFunc("/drivers/{id}", handlers.UpdateDriver).Methods("PUT")
	r.HandleFunc("/drivers/{id}", handlers.DeleteDriver).Methods("DELETE")

	// Rutas para las sesiones
	r.HandleFunc("/sessions", handlers.GetAllSessions).Methods("GET")
	r.HandleFunc("/sessions/{id}", handlers.GetSession).Methods("GET")
	r.HandleFunc("/sessions", handlers.CreateSession).Methods("POST")
	r.HandleFunc("/sessions/{id}", handlers.UpdateSession).Methods("PUT")
	r.HandleFunc("/sessions/{id}", handlers.DeleteSession).Methods("DELETE")

	return r
}
