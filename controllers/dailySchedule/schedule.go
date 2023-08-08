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
// @Description Uploads the daily schedule event
// @Tags Event
// @Accept json
// @Produce json
// @Security Bearer
// @Param schedule body restTypes.Event true "Daily Schedule data to update"
// @Success 201 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [post]
func PostDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data from the request body
	var schedule restTypes.Event
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Validate the schedule data

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
	response := restTypes.LoginResponse{
		Status:  "success",
		Message: "Daily schedule uploaded successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// PutDailySchedule @Summary Update an event in the daily schedule for a given ID and date
// @Description Updates the daily schedule event based on the ID and date provided in the JSON
// @Tags Event
// @Accept json
// @Produce json
// @Security Bearer
// @Param schedule body restTypes.Event  true "Updated Daily Schedule data"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [put]
func PutDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data from the request body
	var schedule restTypes.Event
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract the id and date from the Start time in the JSON data
	id := schedule.ID
	date := schedule.Start.Format("2006-01-02") // Format the date as "YYYY-MM-DD"

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Check if the event exists in the database for the given ID and date
	query := "SELECT COUNT(*) FROM events WHERE id=? AND DATE(start)=?"
	var count int
	err = db.QueryRow(query, id, date).Scan(&count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		// Event not found for the given ID and date
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// Update the event data in the database
	query = "UPDATE events SET title=?, description=?, start=?, end=?, status=?, color=?, location=? WHERE id=? AND DATE(start)=?"
	_, err = db.Exec(query, schedule.Title, schedule.Description, schedule.Start, schedule.End, schedule.Status, schedule.Color, schedule.Location, id, date)
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

// DeleteDailySchedule @Summary Delete an event from the daily schedule for a given ID and date
// @Description Deletes the daily schedule event based on the ID and date provided in the JSON
// @Tags Event
// @Accept json
// @Produce json
// @Security Bearer
// @Param schedule body restTypes.Event  true "Data to delete an event from the daily schedule"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 400 {string} Bad Request
// @Failure 404 {string} Not Found
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/ [delete]
func DeleteDailySchedule(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON data from the request body
	var schedule restTypes.Event
	err := json.NewDecoder(r.Body).Decode(&schedule)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Extract the id and date from the Start time in the JSON data
	id := schedule.ID
	date := schedule.Start.Format("2006-01-02") // Format the date as "YYYY-MM-DD"

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Check if the event exists in the database for the given ID and date
	query := "SELECT COUNT(*) FROM events WHERE id=? AND DATE(start)=?"
	var count int
	err = db.QueryRow(query, id, date).Scan(&count)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		// Event not found for the given ID and date
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// Delete the event from the database
	query = "DELETE FROM events WHERE id=? AND DATE(start)=?"
	_, err = db.Exec(query, id, date)
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

// GetEventsByDate @Summary Get all events for a given date
// @Description Retrieves all events for the specified date
// @Tags Event
// @Accept json
// @Produce json
// @Param date query string true "Date of events to retrieve (in the format 'YYYY-MM-DD')"
// @Success 200 {object} restTypes.GetEventsResponse
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /data/daily-schedule/events [get]
func GetEventsByDate(w http.ResponseWriter, r *http.Request) {
	// Parse the date query parameter from the URL
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse("2006-01-02", dateStr)
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

	// Query all events for the specified date
	query := "SELECT id, title, description, start, end, status, color, location FROM events WHERE DATE(start)=?"
	rows, err := db.Query(query, date.Format("2006-01-02"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	// Collect the events into a slice
	var events []restTypes.Event
	for rows.Next() {
		var event restTypes.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.Start, &event.End, &event.Status, &event.Color, &event.Location)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		events = append(events, event)
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Write the response
	response := restTypes.GetEventsResponse{
		Events: events,
	}
	json.NewEncoder(w).Encode(response)
}

///////////////DAILY IMAGE PART

// GetDailyImage returns the image file for the daily schedule of the specified date or the current date.
//
// Retrieves the image file for the daily schedule of the specified date or the current date from the database.
//
// @Summary Get the image file for the daily schedule of the specified date or the current date
// @Tags Event
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
// @Tags Event
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
// @Tags Event
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
