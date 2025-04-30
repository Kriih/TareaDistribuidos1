package models

type Driver struct {
	DriverNumber int    `json:"driver_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	NameAcronym  string `json:"name_acronym"`
	TeamName     string `json:"team_name"`
	CountryCode  string `json:"country_code"`
}

type DriverDetail struct {
	SessionKey       int     `json:"session_key"`
	CircuitShortName string  `json:"circuit_short_name"`
	Race             string  `json:"race"`
	Position         int     `json:"position"`
	FastestLap       bool    `json:"fastest_lap"`
	MaxSpeed         float64 `json:"max_speed"`
	BestLapDuration  float64 `json:"best_lap_duration"`
}


type PerformanceSummary struct {
	Wins         int     `json:"wins"`
	Top3Finishes int     `json:"top_3_finishes"`
	MaxSpeed     float64 `json:"max_speed"`
}

type DriverDetailConsole struct {
	DriverID           int `json:"driver_id"`
	PerformanceSummary struct {
		Wins        int `json:"wins"`
		Top3Finishes int `json:"top_3_finishes"`
		MaxSpeed    int `json:"max_speed"`
	} `json:"performance_summary"`
	RaceResults []struct {
		SessionKey       int     `json:"session_key"`
		CircuitShortName string  `json:"circuit_short_name"`
		Race             string  `json:"race"`
		Position         int     `json:"position"`
		FastestLap       bool    `json:"fastest_lap"`
		MaxSpeed         int     `json:"max_speed"`
		BestLapDuration  float64 `json:"best_lap_duration"`
	} `json:"race_results"`
}
