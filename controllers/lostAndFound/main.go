package lostAndFound

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/authService"
	"server/databaseTypes"
	"server/restTypes"
	"strconv"
	"strings"
)

// GetLostAndFoundItemsHandler retrieves a list of all lost and found items from the database and returns them as a JSON array in the response body.
// @Summary Get all lost and found items
// @Description Retrieves a list of all lost and found items from the database and returns them as a JSON array in the response body.
// @Tags LostAndFound
// @Accept  json
// @Produce  json
// @Success 200 {array} databaseTypes.LostAndFound
// @Failure 401 {string} Unauthorized
// @Failure 500 {string} Internal Server Error
// @Router /data/lost-and-found/ [get]
func GetLostAndFoundItemsHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute query to get lost and found items
	rows, err := db.Query("SELECT id, item_name, description, date_found, location_found, status FROM LostAndFound")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create response struct
	response := restTypes.LostAndFoundResponse{
		Items: []databaseTypes.LostAndFound{},
	}

	// Parse rows into response struct
	for rows.Next() {
		var item databaseTypes.LostAndFound
		err := rows.Scan(&item.ID, &item.ItemName, &item.Description, &item.DateFound, &item.LocationFound, &item.Status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		item.ImageURL = fmt.Sprintf("/data/lost-and-found/image/%d", item.ID)
		response.Items = append(response.Items, item)
	}

	// Marshal response into JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response headers and body
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}

// GetLostAndFoundImageHandler retrieves the image file for a lost and found item with the specified ID from the database
// and returns it as a JPEG image in the response body.
//
// @Summary Get the image file for a lost and found item by ID.
// @Description Fetches the image file for a lost and found item with the specified ID from the database and returns it as a JPEG image in the response body.
// @Tags LostAndFound
// @Accept  json
// @Produce  image/jpeg
// @Param imageID path string true "The ID of the lost and found item to retrieve the image file for."
// @Success 200 {string} binary "The image file for the specified lost and found item."
// @Failure 404 {string} Not Found "The specified lost and found item ID was not found in the database."
// @Failure 500 {string} Internal Server Error
// @Router /data/lost-and-found/image/{imageID} [get]
func GetLostAndFoundImageHandler(w http.ResponseWriter, r *http.Request) {
	imageID := strings.TrimPrefix(r.URL.Path, "/data/lost-and-found/image/")

	// Fetch image from database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "SELECT image_file FROM LostAndFound WHERE id = ?"
	row := db.QueryRow(query, imageID)

	var image []byte
	err = row.Scan(&image)
	if err != nil {
		log.Fatal(err)
	}

	// Set response headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(image)))

	// Write image data to response
	_, err = w.Write(image)
	if err != nil {
		log.Fatal(err)
	}
}

// PostLostAndFoundItem Add a lost and found item
// @Summary Add a lost and found item
// @Description Adds a lost and found item to the database and uploads an image file
// @Tags LostAndFound
// @Accept multipart/form-data
// @Security Bearer
// @Produce  json
// @Param item_name formData string true "Name of the lost/found item"
// @Param description formData string false "Description of the lost/found item"
// @Param date_found formData string true "Date the item was found"
// @Param location_found formData string true "Location where the item was found"
// @Param status formData string true "Status of the item (lost or found)"
// @Param image_file formData file true "Image of the lost/found item"
// @Success 201 {object} restTypes.LostAndFoundPostResponse
// @Failure 400 {object} restTypes.LostAndFoundErrorResponse
// @Failure 401 {string} Unauthorized
// @Failure 500 {object} restTypes.LostAndFoundErrorResponse
// @Router /data/lost-and-found/ [post]
func PostLostAndFoundItem(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	err := r.ParseMultipartForm(30 << 20) // 30 MB max file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to parse form data"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Parse JSON input
	var lostAndFound restTypes.LostAndFoundInput
	lostAndFound.ItemName = r.FormValue("item_name")
	if lostAndFound.ItemName == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Item name is required"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	lostAndFound.Description = r.FormValue("description")
	lostAndFound.DateFound = r.FormValue("date_found")
	if lostAndFound.DateFound == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Date found is required"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	lostAndFound.LocationFound = r.FormValue("location_found")
	if lostAndFound.LocationFound == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Location found is required"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	lostAndFound.Status = r.FormValue("status")
	if lostAndFound.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Status is required"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	file, _, err := r.FormFile("image_file")
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

	// Insert lost and found item into database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to open database connection"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO LostAndFound (item_name, description, date_found, location_found, status, image_file, submitter_id) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to prepare statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer stmt.Close()
	user, _ := authService.IsAuthorized(w, r)
	submitterID := user.ID
	result, err := stmt.Exec(lostAndFound.ItemName, lostAndFound.Description, lostAndFound.DateFound, lostAndFound.LocationFound, lostAndFound.Status, image, submitterID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to execute statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to get rows affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return

	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "No rows were affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Get ID of the newly inserted row
	id, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to get last insert ID"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Return response
	response := restTypes.LostAndFoundPostResponse{
		Status:  "success",
		Message: "lost and found item added",
		ID:      id,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// PutLostAndFoundItem Update a lost and found item
// @Summary Update a lost and found item
// @Description Update an existing lost and found item in the database with the specified ID
// @Tags LostAndFound
// @ID update-lost-and-found-item
// @Accept json
// @Security Bearer
// @Produce json
// @Param id path int true "Lost and found item ID"
// @Param item_name formData string false "Item name"
// @Param description formData string false "Item description"
// @Param date_found formData string false "Date the item was found"
// @Param location_found formData string false "Location where the item was found"
// @Param status formData string false "Status of the lost and found item"
// @Param image_file formData file false "Image file of the lost and found item"
// @Success 200 {object} restTypes.LostAndFoundErrorResponse
// @Failure 400 {object} restTypes.LostAndFoundErrorResponse
// @Failure 500 {object} restTypes.LostAndFoundErrorResponse
// @Router /data/lost-and-found/image/{id} [put]
func PutLostAndFoundItem(w http.ResponseWriter, r *http.Request) {

	// Parse form data
	err := r.ParseMultipartForm(30 << 20) // 30 MB max file size
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to parse form data"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Parse JSON input
	var itemName sql.NullString
	itemName.String = r.FormValue("item_name")
	if itemName.String != "" {
		itemName.Valid = true
	}
	var description sql.NullString
	description.String = r.FormValue("description")
	if description.String != "" {
		description.Valid = true
	}
	var dateFound sql.NullString
	dateFound.String = r.FormValue("date_found")
	if dateFound.String != "" {
		dateFound.Valid = true
	}
	var locationFound sql.NullString
	locationFound.String = r.FormValue("location_found")
	if locationFound.String != "" {
		locationFound.Valid = true
	}
	var status sql.NullString
	status.String = r.FormValue("status")
	if status.String != "" {
		status.Valid = true
	}
	file, _, err := r.FormFile("image_file")
	var img []byte
	if err == nil {
		defer file.Close()

		// Read the image data from the file
		img, err = ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

	}

	// Get the ID of the item to be updated from the URL parameter
	idStr := strings.TrimPrefix(r.URL.Path, "/data/lost-and-found/image/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Invalid ID"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Get the current values of the lost and found item from the database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to open database connection"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT item_name, description, date_found, location_found, status, image_file FROM LostAndFound WHERE id=?", id)
	var currentItemName, currentDescription, currentDateFound, currentLocationFound, currentStatus string
	var currentImage []byte
	err = row.Scan(&currentItemName, &currentDescription, &currentDateFound, &currentLocationFound, &currentStatus, &currentImage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to get current values of lost and found item"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Update the lost and found item in the database
	stmt, err := db.Prepare("UPDATE LostAndFound SET item_name=COALESCE(?,?), description=COALESCE(?,?), date_found=COALESCE(?,?), location_found=COALESCE(?,?), status=COALESCE(?,?), image_file=COALESCE(?,?) WHERE id=?")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to prepare statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(itemName, currentItemName, description, currentDescription, dateFound, currentDateFound, locationFound, currentLocationFound, status, currentStatus, img, currentImage, id)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to execute statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "Failed to get rows affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.LostAndFoundErrorResponse{Error: "No rows were affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Return response
	response := restTypes.LostAndFoundPostResponse{
		Status:  "success",
		Message: "lost and found item updated",
		ID:      int64(id),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type deleteResponse struct {
	Status string `json:"status"`
}

// HandleDeleteLostAndFound  Delete a lost and found item
// @Summary Delete a lost and found item
// @Description Deletes a lost and found item from the database
// @Tags LostAndFound
// @ID delete-lost-and-found-item
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} deleteResponse
// @Failure 400 {string} string "Invalid item ID"
// @Failure 404 {string} string "Item not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Router /lost-and-found/{id} [delete]
func HandleDeleteLostAndFound(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the item_id from the URL path
	urlPath := strings.Split(r.URL.Path, "/")
	itemID := urlPath[len(urlPath)-1]

	// Check if the item_id is valid
	if _, err := strconv.Atoi(itemID); err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete the item from the LostAndFound table
	result, err := db.Exec("DELETE FROM LostAndFound WHERE id=?", itemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the item was deleted successfully
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Create the response struct
	response := deleteResponse{Status: "Item deleted successfully"}

	// Return a JSON response with the status of the operation
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
