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
			fmt.Printf("| %-10s | %-12s | %-15s | %-4s |\n", "Nombre", "Apellido", "Equipo", "País")
			fmt.Println("-------------------------------------------------------------")
			for _, driver := range drivers {
				fmt.Printf("| %-10s | %-12s | %-15s | %-4s |\n", driver.FirstName, driver.LastName, driver.TeamName, driver.CountryCode)
			}

		case "2":
			num := readInput("Ingrese el numero del piloto: ")
			fmt.Printf("-> Ver detalle del corredor %s\n", num)
			// TODO: GET http://localhost:8080/api/corredor/detalle/{num}
			resp, err := http.Get("http://localhost:8080/api/corredor/detalle/" + num)
			if err != nil {
				fmt.Println("Error al hacer la solicitud:", err)
				return
			}
			defer resp.Body.Close()

			var driverDetail []models.DriverDetail
			if err := json.NewDecoder(resp.Body).Decode(&driverDetail); err != nil {
				fmt.Println("Error al decodificar la respuesta:", err)
				return
			}
			fmt.Println("Detalle del corredor:")
			fmt.Println("-------------------------------------------------------------")
			fmt.Printf("| %-10s | %-12s | %-15s | %-4s |\n", "Nombre", "Apellido", "Carrera", "Posicion")
			fmt.Println("-------------------------------------------------------------")
			for _, detail := range driverDetail {
				//imprime todo lo que consigue sin nada más escrito
				fmt.Printf(detail)
			}

		case "3":
			fmt.Println("-> Ver carreras")
			// TODO: GET http://localhost:8080/api/carrera
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
