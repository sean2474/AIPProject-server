package dailySchedule

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"server/restTypes"
	"strings"
	"time"
)

// GetDailySchedule returns the HTML template for the daily schedule of the specified date or the current date.
//
// Retrieves the daily schedule HTML template for the specified date or the current date from the database.
//
// @Summary Get the daily schedule HTML template for the specified date or the current date
// @Tags DailySchedule
// @Accept  */*
// @Produce  text/html
// @Param date query string false "The date for which to retrieve the daily schedule in the format 'YYYY-MM-DD'. If not provided, the current date is used."
// @Success 200 {string} OK
// @Failure 401 {string} Unauthorized
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [get]
func GetDailySchedule(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get the date parameter from the request, or use the current date if it's not provided
	dateParam := r.URL.Query().Get("date")
	date := time.Now().Format("2006-01-02")
	if dateParam != "" {
		date = dateParam
	}

	// Query the database for the daily schedule HTML template for the specified date
	query := "SELECT html_page FROM DailySchedule WHERE date = ?"
	row := db.QueryRow(query, date)

	// Extract the HTML template data from the row
	var htmlTemplate string
	err = row.Scan(&htmlTemplate)
	if err == sql.ErrNoRows {
		fmt.Println(err.Error())
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the HTML template
	tmpl, err := template.New("dailySchedule").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with no data and write the output to the response body
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PostDailySchedule @Summary Upload the daily schedule image for the current date
// @Summary Upload the daily schedule image for the current date
// @Description Uploads the daily schedule image for the current date to the database
// @Tags DailySchedule
// @Accept  multipart/form-data
// @Produce  json
// @Security Bearer
// @Param image formData file true "The daily schedule image file"
// @Success 201 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [post]
func PostDailySchedule(w http.ResponseWriter, r *http.Request) {

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
	image, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the image data into the database
	query := "INSERT INTO DailySchedule (date, image_file, image_url, file_name) VALUES (?, ?, ?, ?)"
	_, err = db.Exec(query, time.Now().Format("2006-01-02"), image, "", "")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	// Write the response
	response := restTypes.LoginResponse{
		Status:  "success",
		Message: "Daily schedule uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// PutDailySchedule @Summary Update the daily schedule image for a specific date
// @Summary Update the daily schedule for a specific date
// @Description Updates the daily schedule image for a specific date in the database
// @Tags DailySchedule
// @Produce  json
// @Security Bearer
// @Param date path string true "The date of the daily schedule to update (in YYYY-MM-DD format)"
// @Param image formData file true "The updated image file for the daily schedule (limit: 32 MB)"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/{date} [put]
func PutDailySchedule(w http.ResponseWriter, r *http.Request) {

	// Get the schedule date from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	scheduleDate := parts[3]

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
	image, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Update the image data in the database using the date
	query := "UPDATE DailySchedule SET image_file = ? WHERE date = ?"
	_, err = db.Exec(query, image, scheduleDate)
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
		Message: "Daily schedule updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteDailySchedule @Summary Delete the daily schedule image for a specific date
// @Summary Delete the daily schedule for a specific date
// @Description Deletes the daily schedule for a specific date from the database
// @Tags DailySchedule
// @Produce  json
// @Security Bearer
// @Param date path string true "The date of the daily schedule to delete (in YYYY-MM-DD format)"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/{date} [delete]
func DeleteDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Get the schedule date from the URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	scheduleDate := parts[3]
	fmt.Println(scheduleDate)
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Delete the record from the database using the date
	query := "DELETE FROM DailySchedule WHERE date = ?"
	_, err = db.Exec(query, scheduleDate)
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
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
