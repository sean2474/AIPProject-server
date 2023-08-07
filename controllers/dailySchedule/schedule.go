package dailySchedule

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/restTypes"
	"time"
)

// PostDailySchedule @Summary Upload the daily schedule image for the current date
// @Description Uploads the daily schedule image for the current date to the database
// @Tags DailySchedule
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param image formData file true "The daily schedule image file"
// @Success 201 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [post]
func PostDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data from the request body
	var schedule restTypes.DailySchedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Insert the schedule data into the database
	query := "INSERT INTO events (id, title, description, start, end, status, color, location) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Exec(query, schedule.ID, schedule.Title, schedule.Description, schedule.Start, schedule.End, schedule.Status, schedule.Color, schedule.Location)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// Write the response
	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Daily schedule uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// PutDailySchedule @Summary Update the daily schedule
// @Description Updates the daily schedule with the provided data
// @Tags DailySchedule
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Daily Schedule ID to update"
// @Param schedule body DailySchedule true "Daily Schedule data to update"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/{id} [put]
func PutDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Function body for PUT request
}

// DeleteDailySchedule @Summary Delete the daily schedule
// @Description Deletes the daily schedule with the provided ID
// @Tags DailySchedule
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Daily Schedule ID to delete"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/{id} [delete]
func DeleteDailySchedule(w http.ResponseWriter, r *http.Request) {
	var schedule restTypes.DailySchedule
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Delete the schedule data from the database
	query := "DELETE FROM DailySchedule WHERE id=?"
	_, err = db.Exec(query, schedule.ID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	response := restTypes.DeleteResponse{
		Status:  "success",
		Message: "Daily schedule deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// GetDailyImage returns the image file for the daily schedule of the specified date or the current date.
//
// Retrieves the image file for the daily schedule of the specified date or the current date from the database.
//
// @Summary Get the image file for the daily schedule of the specified date or the current date
// @Tags DailySchedule
// @Accept  */*
// @Produce  image/*
// @Param date query string false "The date for which to retrieve the daily schedule image in the format 'YYYY-MM-DD'. If not provided, the current date is used."
// @Success 200 {string} OK
// @Failure 401 {string} Unauthorized
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/image [get]
func GetDailyImage(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the date parameter from the request, or use the current date if it's not provided
	dateParam := r.URL.Query().Get("date")
	date := time.Now().Format("2006-01-02")
	if dateParam != "" {
		// Attempt to parse the provided date
		t, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			http.Error(w, "Invalid date format. Please use 'YYYY-MM-DD' format.", http.StatusBadRequest)
			return
		}
		// Use the parsed date
		date = t.Format("2006-01-02")
	}

	// Query the database for the daily schedule image for the specified date
	query := "SELECT image_file FROM DailyScheduleImages WHERE date = ?"
	row := db.QueryRow(query, date)

	// Extract the image data from the row
	var imageData []byte
	err = row.Scan(&imageData)
	if err == sql.ErrNoRows {
		// If no record is found in the database, serve the "404.jpg" image instead
		imageData, err = ioutil.ReadFile("static/404.png")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Set the response headers and write the image data to the response body
	w.Header().Set("Content-Type", "image/*")
	w.WriteHeader(http.StatusOK)
	w.Write(imageData)
}

// PostDailyImage @Summary Upload the daily schedule image for a specific date
// @Description Uploads the daily schedule image for the provided date to the database
// @Tags DailySchedule
// @Accept  multipart/form-data
// @Produce  json
// @Security Bearer
// @Param image formData file true "The daily schedule image file"
// @Param date formData string true "The date for which the image is uploaded (format: 2006-01-02)"
// @Success 201 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/image [post]
func PostDailyImage(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseMultipartForm(32 << 20) // Limit: 32 MB
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get the image file from the form data
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the image data from the file
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Get the date parameter from the form data
	date := r.FormValue("date")
	if date == "" {
		http.Error(w, "Bad Request - 'date' parameter is required", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if there is an existing record for the provided date
	query := "SELECT id FROM DailyScheduleImages WHERE date = ?"
	row := db.QueryRow(query, date)

	var existingID int
	err = row.Scan(&existingID)
	if err == nil {
		// If an existing record is found, update the image data
		updateQuery := "UPDATE DailyScheduleImages SET image_file = ? WHERE date = ?"
		_, err = db.Exec(updateQuery, imageData, date)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if err == sql.ErrNoRows {
		// If no record is found for the provided date, insert a new one
		insertQuery := "INSERT INTO DailyScheduleImages (date, image_file, image_url, file_name) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, date, imageData, "", "")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Error occurred while querying the database
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// Write the response
	response := restTypes.LoginResponse{
		Status:  "success",
		Message: "Daily schedule image uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteDailyImage @Summary Delete the daily schedule image for a specific date
// @Description Deletes the daily schedule image for a specific date from the database
// @Tags DailySchedule
// @Produce  json
// @Security Bearer
// @Param date query string false "The date for which to delete the daily schedule image in the format 'YYYY-MM-DD'. If not provided, the current date is used."
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/image [delete]
func DeleteDailyImage(w http.ResponseWriter, r *http.Request) {
	// Get the date parameter from the request, or use the current date if it's not provided
	dateParam := r.URL.Query().Get("date")
	date := time.Now().Format("2006-01-02")
	if dateParam != "" {
		date = dateParam
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete the record from the database using the date
	query := "DELETE FROM DailyScheduleImages WHERE date = ?"
	_, err = db.Exec(query, date)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	response := restTypes.LoginResponse{
		Status:  "success",
		Message: "Daily schedule image deleted successfully",
	}
	json.NewEncoder(w).Encode(response)
}
