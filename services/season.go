package services

import (
	"f1-statshub/db"
)

// GetSeasonSummary genera el resumen de la temporada
func GetSeasonSummary() (map[string]interface{}, error) {
	type TopResult struct {
		Position int    `json:"position"`
		Driver   string `json:"driver"`
		Team     string `json:"team"`
		Country  string `json:"country"`
		Value    int    `json:"value"`
	}

	var topWins []TopResult
	var topFastestLaps []TopResult
	var topPoles []TopResult

	// Top 3 pilotos con más victorias
	rows, err := db.DB.Query(`
		SELECT d.first_name || ' ' || d.last_name as driver, d.team_name, d.country_code, COUNT(*) as wins
		FROM positions p
		JOIN drivers d ON p.driver_number = d.driver_number
		WHERE p.position = 1
		GROUP BY d.driver_number
		ORDER BY wins DESC
		LIMIT 3
	`)
	if err == nil {
		defer rows.Close()
		pos := 1
		for rows.Next() {
			var r TopResult
			err := rows.Scan(&r.Driver, &r.Team, &r.Country, &r.Value)
			if err == nil {
				r.Position = pos
				topWins = append(topWins, r)
				pos++
			}
		}
	}

	// Top 3 pilotos con más vueltas rápidas (fastest laps)
	rows, err = db.DB.Query(`
		SELECT d.first_name || ' ' || d.last_name as driver, d.team_name, d.country_code, COUNT(*) as fastest_laps
		FROM laps l
		JOIN drivers d ON l.driver_number = d.driver_number
		WHERE l.lap_number = 1
		GROUP BY d.driver_number
		ORDER BY fastest_laps DESC
		LIMIT 3
	`)
	if err == nil {
		defer rows.Close()
		pos := 1
		for rows.Next() {
			var r TopResult
			err := rows.Scan(&r.Driver, &r.Team, &r.Country, &r.Value)
			if err == nil {
				r.Position = pos
				topFastestLaps = append(topFastestLaps, r)
				pos++
			}
		}
	}

	// Top 3 pilotos con más pole positions
	rows, err = db.DB.Query(`
		SELECT d.first_name || ' ' || d.last_name as driver, d.team_name, d.country_code, COUNT(*) as poles
		FROM positions p
		JOIN drivers d ON p.driver_number = d.driver_number
		WHERE p.position = 1
		GROUP BY d.driver_number
		ORDER BY poles DESC
		LIMIT 3
	`)
	if err == nil {
		defer rows.Close()
		pos := 1
		for rows.Next() {
			var r TopResult
			err := rows.Scan(&r.Driver, &r.Team, &r.Country, &r.Value)
			if err == nil {
				r.Position = pos
				topPoles = append(topPoles, r)
				pos++
			}
		}
	}

	return map[string]interface{}{
		"season":               2024,
		"top_3_winners":        topWins,
		"top_3_fastest_laps":   topFastestLaps,
		"top_3_pole_positions": topPoles,
	}, nil
}
