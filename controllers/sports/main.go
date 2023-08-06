package sports

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/databaseTypes"
	"server/restTypes"
	"time"
)

// GetSportsData Get sports data
// @Summary Get sports data
// @Description Retrieves data about sports teams and their coaches.
// @ID get-sports-data
// @Tags SportsData
// @Accept json
// @Produce json
// @Success 200 {object} restTypes.SportsDataList "List of sports data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /data/sports/ [get]
func GetSportsData(w http.ResponseWriter, r *http.Request) {

	// Connect to the database
	db, err := sql.Open("sqlite3", "database.db")
	fmt.Println("1")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	fmt.Println("2")
	// Query the databaseTypes.SportsInfo table to get the data
	rows, err := db.Query("SELECT id, sport_name, category, season, coach_name, coach_contact, roster FROM SportsInfo")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice of databaseTypes.SportsInfo structs to hold the data
	var sportsDataList restTypes.SportsDataList

	// Loop through the rows and add the data to the slice
	for rows.Next() {
		var sportsData databaseTypes.SportsInfo
		err := rows.Scan(&sportsData.ID, &sportsData.SportName, &sportsData.Category, &sportsData.Season,
			&sportsData.CoachName, &sportsData.CoachContact, &sportsData.Roster)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Add the databaseTypes.SportsInfo struct to the slice
		sportsDataList.List = append(sportsDataList.List, sportsData)
	}

	// Marshal the slice to JSON
	jsonData, err := json.Marshal(sportsDataList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON data to the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetSportsGameData Get sports game data
// @Summary Get sports game data
// @Description Retrieves data about sports games and their results.
// @ID get-sports-game-data
// @Tags SportsData
// @Accept json
// @Produce json
// @Success 200 {object} restTypes.SportsGameDataList "List of sports game data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /data/games/ [get]
func GetSportsGameData(w http.ResponseWriter, r *http.Request) {

	// Connect to the database
	db, err := sql.Open("sqlite3", "database.db?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the SportsGames table to get the data
	rows, err := db.Query("SELECT SportsGames.id, SportsGames.sport_name, SportsGames.category, SportsGames.game_location, SportsGames.opponent_school, SportsGames.home_or_away, SportsGames.match_result, SportsGames.coach_comment, strftime(SportsGames.game_schedule) FROM SportsGames")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice of sportsGameData structs to hold the data
	var sportsGameDataList restTypes.SportsGameDataList

	// Loop through the rows and add the data to the slice
	for rows.Next() {
		var sportsGameData databaseTypes.SportsGame
		var gameSchedule string
		err := rows.Scan(&sportsGameData.ID, &sportsGameData.SportName, &sportsGameData.Category, &sportsGameData.GameLocation,
			&sportsGameData.OpponentSchool, &sportsGameData.HomeOrAway, &sportsGameData.MatchResult, &sportsGameData.CoachComment,
			&gameSchedule)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(gameSchedule)
		layout := "2006-01-02 03:04 PM"

		t, err := time.Parse(layout, gameSchedule)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return
		}
		// Format the game schedule as a string
		sportsGameData.GameSchedule = t

		// Add the sportsGameData struct to the slice
		sportsGameDataList.List = append(sportsGameDataList.List, sportsGameData)
	}

	// Marshal the slice to JSON
	jsonData, err := json.Marshal(sportsGameDataList)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set response headers and write the JSON data to the response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
