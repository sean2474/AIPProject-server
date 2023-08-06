package food

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/databaseTypes"
	"server/restTypes"
	"strings"
	"time"
)

// PostFoodMenuHandler @Summary Add a food menu
// @Summary Add a food menu
// @Description Add a new food menu to the database
// @Tags FoodMenu
// @Security Bearer
// @ID addFoodMenu
// @Accept json
// @Produce json
// @Param foodMenu body databaseTypes.FoodMenu true "Food menu to add"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /data/food-menu/ [post]
func PostFoodMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Verify authentication here if needed

	var foodMenu databaseTypes.FoodMenu
	err := json.NewDecoder(r.Body).Decode(&foodMenu)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO FoodMenu (date, breakfast, lunch, dinner) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(foodMenu.Date, foodMenu.Breakfast, foodMenu.Lunch, foodMenu.Dinner)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restTypes.LoginResponse{
		Status:  "success",
		Message: "Food menu added successfully",
	})
}

// PutFoodMenuHandler handles a PUT request to /food-menu/{id}
// @Summary Update a food menu
// @Description Update the food menu with the specified ID
// @Tags FoodMenu
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path string true "ID of the food menu to update"
// @Param foodMenu body databaseTypes.FoodMenu true "New values for the food menu"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {object} restTypes.ErrorResponse
// @Failure 401 {object} restTypes.ErrorResponse
// @Failure 404 {object} restTypes.ErrorResponse
// @Router /data/food-menu/{id} [put]
func PutFoodMenuHandler(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//TODO: Verify authentication here if needed

	var foodMenu databaseTypes.FoodMenu
	err := json.NewDecoder(r.Body).Decode(&foodMenu)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE FoodMenu SET date=?, breakfast=?, lunch=?, dinner=? WHERE date=?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(foodMenu.Date, foodMenu.Breakfast, foodMenu.Lunch, foodMenu.Dinner, id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(restTypes.LoginResponse{
		Status:  "success",
		Message: "Food menu updated successfully",
	})
}

// DeleteFoodMenu @Summary Delete a food menu
// @Summary Delete a food menu
// @Description Delete a food menu from the database for a given date
// @Tags FoodMenu
// @Security Bearer
// @ID DeleteFoodMenu
// @Param date path string true "The date of the food menu to delete"
// @Success 200 {object} restTypes.DeleteResponse
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /data/food-menu/{date} [delete]
func DeleteFoodMenu(w http.ResponseWriter, r *http.Request, date string) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete the food menu for the given date
	stmt, err := db.Prepare("DELETE FROM FoodMenu WHERE date = ?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := stmt.Exec(date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if any row was affected by the delete operation
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Return a success message in JSON format
	response := restTypes.DeleteResponse{Status: "success", Message: "Food menu deleted successfully"}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// GetFoodMenu @Summary Get the food menu for the current date
// @Summary Get the food menu for the current date
// @Description Retrieves the breakfast, lunch, and dinner menu for the current date from the database
// @Tags FoodMenu
// @Accept  json
// @Produce  json
// @Success 200 {object} databaseTypes.FoodMenu
// @Failure 401 {string} Unauthorized
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/food-menu/ [get]
func GetFoodMenu(w http.ResponseWriter, r *http.Request) {

	// Get the current date
	date := time.Now()
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Query the database for the food menu for the current date
	query := "SELECT breakfast, lunch, dinner FROM FoodMenu WHERE date = ?"
	row := db.QueryRow(query, date.Format("2006-01-02"))
	fmt.Println(date.Format("2006-01-02"))
	// Extract the values from the row
	var breakfast, lunch, dinner string
	err = row.Scan(&breakfast, &lunch, &dinner)
	if err == sql.ErrNoRows {
		fmt.Println(err.Error())
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a FoodMenu struct from the values
	foodMenu := databaseTypes.FoodMenu{
		Date:      date.Format("2006-01-02"),
		Breakfast: breakfast,
		Lunch:     lunch,
		Dinner:    dinner,
	}

	// Convert the FoodMenu struct to a JSON object
	jsonData, err := json.Marshal(foodMenu)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	w.Write(jsonData)
}

// GetFoodMenuByDate @Summary	 Get the food menu for a specific date
// @Description Retrieves the breakfast, lunch, and dinner menu for a specific date from the database
// @Tags FoodMenu
// @Accept  json
// @Produce  json
// @Param   date      path    string    true        "The date of the food menu (YYYY-MM-DD)"
// @Success 200 {object} databaseTypes.FoodMenu
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/food-menu/{date} [get]
func GetFoodMenuByDate(w http.ResponseWriter, r *http.Request) {
	// Get the date parameter from the path
	dateStr := strings.TrimPrefix(r.URL.Path, "/data/food-menu/")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the database for the food menu for the specified date
	query := "SELECT breakfast, lunch, dinner FROM FoodMenu WHERE date = ?"
	row := db.QueryRow(query, date.Format("2006-01-02"))

	// Extract the values from the row
	var breakfast, lunch, dinner string
	err = row.Scan(&breakfast, &lunch, &dinner)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a FoodMenu struct from the values
	foodMenu := databaseTypes.FoodMenu{
		Date:      date.Format("2006-01-02"),
		Breakfast: breakfast,
		Lunch:     lunch,
		Dinner:    dinner,
	}

	// Convert the FoodMenu struct to a JSON object
	jsonData, err := json.Marshal(foodMenu)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	w.Write(jsonData)
}

// GetAllFoodMenus @Summary Get all the food menus from the database
// @Description Retrieves all the breakfast, lunch, and dinner menus from the database
// @Tags FoodMenu
// @Accept json
// @Produce json
// @Success 200 {object} restTypes.AllMenuResponse
// @Failure 500 {string} Internal Server Error
// @Router /data/food-menu/all [get]
func GetAllFoodMenus(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the database to get all food menus
	query := "SELECT date, breakfast, lunch, dinner FROM FoodMenu"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to store all food menus
	var foodMenus restTypes.AllMenuResponse

	// Iterate over the rows and extract the values
	for rows.Next() {
		var date, breakfast, lunch, dinner string
		err := rows.Scan(&date, &breakfast, &lunch, &dinner)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Create a FoodMenu struct for each row
		foodMenu := databaseTypes.FoodMenu{
			Date:      date,
			Breakfast: breakfast,
			Lunch:     lunch,
			Dinner:    dinner,
		}
		foodMenus.Items = append(foodMenus.Items, foodMenu)
		// Append the food menu to the slice
		//foodMenus = append(foodMenus, foodMenu)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the foodMenus slice to a JSON object
	jsonData, err := json.Marshal(foodMenus)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	w.Write(jsonData)
}
