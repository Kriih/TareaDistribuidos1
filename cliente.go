package main

import (
	"bufio"
	"encoding/json"
	"f1-statshub/models"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	for {
		showMenu()
		option := readInput("Seleccione una opcion: ")

		switch option {
		case "1":
			fmt.Println("-> Ver corredores")
			// TODO: GET http://localhost:8080/api/corredor
			// hacer el get a la api y mostrar los resultados
			resp, err := http.Get("http://localhost:8080/api/corredor")
			if err != nil {
				fmt.Println("Error al hacer la solicitud:", err)
				return
			}
			defer resp.Body.Close()

			var drivers []models.Driver
			if err := json.NewDecoder(resp.Body).Decode(&drivers); err != nil {
				fmt.Println("Error al decodificar la respuesta:", err)
				return
			}
			fmt.Println("Lista de corredores:")
			fmt.Println("-------------------------------------------------------------")
			fmt.Printf("| %-5s | %-10s | %-12s | %-15s | %-4s |\n", "Nro", "Nombre", "Apellido", "Equipo", "País")
			fmt.Println("-------------------------------------------------------------")
			for _, driver := range drivers {
				fmt.Printf("| %-5d | %-10s | %-12s | %-15s | %-4s |\n", driver.DriverNumber, driver.FirstName, driver.LastName, driver.TeamName, driver.CountryCode)
			}

		case "2":
			num := readInput("Ingrese el numero del piloto: ")
			fmt.Printf("-> Ver detalle del corredor %s\n", num)
		
			resp, err := http.Get("http://localhost:8080/api/corredor/detalle/" + num)
			if err != nil {
				fmt.Println("Error al hacer la solicitud:", err)
				return
			}
			defer resp.Body.Close()
		
			var driverDetail models.DriverDetailConsole
			if err := json.NewDecoder(resp.Body).Decode(&driverDetail); err != nil {
				fmt.Println("Error al decodificar la respuesta:", err)
				return
			}
		
			// Mostrar resumen de rendimiento
			fmt.Println("\n===== RESUMEN DE RENDIMIENTO =====")
			fmt.Printf("Victorias: %d\n", driverDetail.PerformanceSummary.Wins)
			fmt.Printf("Top 3:     %d\n", driverDetail.PerformanceSummary.Top3Finishes)
			fmt.Printf("Vel. Máx.: %d km/h\n", driverDetail.PerformanceSummary.MaxSpeed)
		
			// Mostrar resultados de carrera
			fmt.Println("\n========= RESULTADOS DE CARRERA =========")
			fmt.Printf("| %-15s | %-25s | %-8s | %-12s | %-10s | %-6s |\n", "Circuito", "Carrera", "Posición", "Vuelta rápida", "Vel. máx", "Mejor vuelta")
			fmt.Println(strings.Repeat("-", 95))
			for _, race := range driverDetail.RaceResults {
				fmt.Printf("| %-15s | %-25s | %-8d | %-12t | %-10d | %-6.3f |\n",
					race.CircuitShortName, race.Race, race.Position, race.FastestLap, race.MaxSpeed, race.BestLapDuration)
			}
		
		

		case "3":
			fmt.Println("-> Ver carreras")
			// TODO: GET http://localhost:8080/api/carrera

			resp, err := http.Get("http://localhost:8080/api/carrera")
			if err != nil {
				fmt.Println("Error al hacer la solicitud:", err)
				return
			}
			defer resp.Body.Close()

			var session []models.Session
			if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
				fmt.Println("Error al decodificar la respuesta:", err)
				return
			}
			fmt.Println("Lista de carreras:")
			fmt.Println("-------------------------------------------------------------") // ID carrera, Pais, Fecha, Year, Circuito
			// type Session struct {
			// 	SessionKey       int    `json:"session_key"`
			// 	SessionName      string `json:"session_name"`
			// 	SessionType      string `json:"session_type"`
			// 	Location         string `json:"location"`
			// 	CountryName      string `json:"country_name"`
			// 	Year             int    `json:"year"`
			// 	CircuitShortName string `json:"circuit_short_name"`
			// 	DateStart        string `json:"date_start"` // ISO8601 string
			// }
			fmt.Printf("| %-10s | %-12s | %-15s | %-4s |\n", "ID Carrera", "Pais", "Fecha", "Circuito")
			fmt.Println("-------------------------------------------------------------")
			for _, session := range session {
				fmt.Printf("| %-5d | %-20s | %-15d | %-4s |\n", session.SessionKey, session.CountryName, session.Year, session.CircuitShortName)
			}

		case "4":
			num := readInput("Ingrese el ID de la carrera: ")
			fmt.Printf("-> Ver detalle de la carrera %s\n", num)
			// TODO: GET http://localhost:8080/api/carrera/detalle/{num}
		case "5":
			year := readInput("Ingrese la temporada: ")
			fmt.Printf("-> Ver resumen de temporada %s\n", year)
			// TODO: GET http://localhost:8080/api/temporada/resumen/
		case "6":
			fmt.Println("Fin del programa!")
			return
		default:
			fmt.Println("Opcion invalida. Intente de nuevo.")
		}
	}
}

func showMenu() {
	fmt.Println("\nMenu")
	fmt.Println("1. Ver corredores")
	fmt.Println("2. Ver detalle de corredor")
	fmt.Println("3. Ver carreras")
	fmt.Println("4. Ver detalle de carrera")
	fmt.Println("5. Resumen de temporada")
	fmt.Println("6. Salir")
}

func readInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
