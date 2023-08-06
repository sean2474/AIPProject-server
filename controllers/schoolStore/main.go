package schoolStore

import (
	"database/sql"
	"encoding/json"
	_ "errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"server/databaseTypes"
	"server/restTypes"
	"strconv"
	"strings"
)

// HandleSchoolStore Get a list of items from the School Store
// @Summary Get a list of items from the School Store
// @Description Retrieves a list of items from the School Store database
// @Tags School Store
// @Accept  json
// @Produce  json
// @Success 200 {object} restTypes.SchoolStoreResponse
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /data/school-store/ [get]
func HandleSchoolStore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT Product_Name, Description, Price, Category, ID FROM School_Store")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []databaseTypes.SchoolStore
	for rows.Next() {
		item := databaseTypes.SchoolStore{}
		err := rows.Scan(&item.ProductName, &item.Description, &item.Price, &item.Category, &item.ID)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	response := restTypes.SchoolStoreResponse{
		List: items,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// HandleSchoolStoreImage Get an image for an item from the School Store
// @Summary Get an image for an item from the School Store
// @Description Retrieves an image for a specified item from the School Store database
// @Tags School Store
// @Accept  json
// @Produce  jpeg
// @Param   item_id    path    int     true        "ID of the item to retrieve the image for"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Internal server error"
// @Router /data/school-store/image/{item_id} [get]
func HandleSchoolStoreImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get item ID from URL parameter
	itemIDStr := strings.TrimPrefix(r.URL.Path, "/data/school-store/image/")
	if itemIDStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var image []byte
	err = db.QueryRow("SELECT image_file FROM School_Store WHERE ID = ?", itemID).Scan(&image)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

// HandleAddSchoolStoreItem Add an item to the School Store
// @Summary Add an item to the School Store
// @Description Adds a new item to the School Store database
// @Tags School Store
// @Security Bearer
// @Accept  multipart/form-data
// @Produce  json
// @Param   item_name    formData    string     true        "Name of the item to add"
// @Param   description  formData    string     true        "Description of the item to add"
// @Param   price        formData    number     true        "Price of the item to add"
// @Param   category     formData    string     true        "Category of the item to add"
// @Param   image_file   formData    file       true        "Image file of the item to add"
// @Success 200 {object} restTypes.ErrorResponse
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 500 {string} string "Internal server error"
// @Router /data/school-store/ [post]
func HandleAddSchoolStoreItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // Max file size: 32 MB
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Extract form values
	itemName := r.FormValue("item_name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	category := r.FormValue("category")
	imageFile, _, err := r.FormFile("image_file")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer imageFile.Close()

	// Convert price to float64
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Read image file into byte slice
	imageBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Insert new item into database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO School_Store (Product_Name, Description, Price, Category, image_file, Date_Added, Stock) VALUES (?, ?, ?, ?, ?, DATE('now'),1)", itemName, description, price, category, imageBytes)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	response := restTypes.ErrorResponse{
		Code:    200,
		Message: "Item added to School Store",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// PutSchoolStoreItem Update an item in the School Store
// @Summary Update an item in the School Store
// @Security Bearer
// @Description Updates an existing item in the School Store database
// @Tags School Store
// @Accept  multipart/form-data
// @Produce  json
// @Param   item_id      path        int         true        "ID of the item to update"
// @Param   item_name    formData    string      false       "New name of the item"
// @Param   description  formData    string      false       "New description of the item"
// @Param   price        formData    number      false       "New price of the item"
// @Param   category     formData    string      false       "New category of the item"
// @Param   image_file   formData    file        false       "New image file of the item"
// @Success 200 {object} restTypes.SchoolStorePostResponse
// @Failure 400 {string} restTypes.SchoolStoreErrorResponse
// @Failure 404 {string} restTypes.SchoolStoreErrorResponse
// @Failure 500 {string} restTypes.SchoolStoreErrorResponse
// @Router /data/school-store/{item_id} [put]
func PutSchoolStoreItem(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseMultipartForm(30 << 20) // 30 MB max file size
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to parse form data"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Parse JSON input
	var productName sql.NullString
	productName.String = r.FormValue("item_name")
	if productName.String != "" {
		productName.Valid = true
	}
	var category sql.NullString
	category.String = r.FormValue("category")
	if category.String != "" {
		category.Valid = true
	}
	var price sql.NullFloat64
	priceStr := r.FormValue("price")
	if priceStr != "" {
		priceVal, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Invalid price"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		price.Float64 = priceVal
		price.Valid = true
	}
	var stock sql.NullInt64
	stockStr := r.FormValue("stock")
	if stockStr != "" {
		stockVal, err := strconv.ParseInt(stockStr, 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Invalid stock"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		stock.Int64 = stockVal
		stock.Valid = true
	}
	var description sql.NullString
	description.String = r.FormValue("description")
	if description.String != "" {
		description.Valid = true
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
	idStr := strings.TrimPrefix(r.URL.Path, "/data/school-store/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Invalid ID"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Get the current values of the school store item from the database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to open database connection"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT product_name, category, price, stock, description, image_file FROM School_Store WHERE id=?", id)
	var currentProductName, currentCategory, currentDescription string
	var currentPrice sql.NullFloat64
	var currentStock sql.NullInt64
	var currentImage []byte
	err = row.Scan(&currentProductName, &currentCategory, &currentPrice, &currentStock, &currentDescription, &currentImage)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to get current values of school store item"}
		json.NewEncoder(w).Encode(errorResponse)
		return

	}
	// Update the school store item in the database
	stmt, err := db.Prepare("UPDATE School_Store SET Product_Name=COALESCE(?,?), Category=COALESCE(?,?), Price=COALESCE(?,?), Stock=COALESCE(?,?), Description=COALESCE(?,?), image_file=COALESCE(?, ?) WHERE ID=?")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to prepare statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(productName, currentProductName, category, currentCategory, price, currentPrice, stock, currentStock, description, currentDescription, img, currentImage, id)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to execute statement"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "Failed to get rows affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := restTypes.SchoolStoreErrorResponse{Error: "No rows were affected"}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Return response
	response := restTypes.SchoolStorePostResponse{
		Status:  "success",
		Message: "school store item updated",
		ID:      int64(id),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// HandleDeleteSchoolStoreItem Delete an item from the School Store
// @Summary Delete an item from the School Store
// @Description Deletes an item from the School Store database
// @Tags School Store
// @Security Bearer
// @Produce  json
// @Param   item_id      path        int         true        "ID of the item to delete"
// @Success 200 {object} restTypes.ErrorResponse
// @Failure 400 {string} string "Invalid request parameters"
// @Failure 404 {string} string "Item not found"
// @Failure 500 {string} string "Internal server error"
// @Router /data/school-store/{item_id} [delete]
func HandleDeleteSchoolStoreItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get item ID from URL parameter
	itemIDStr := strings.TrimPrefix(r.URL.Path, "/data/school-store/")
	if itemIDStr == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Delete item from database
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var image []byte
	err = db.QueryRow("SELECT image_file FROM School_Store WHERE ID = ?", itemID).Scan(&image)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	result, err := db.Exec("DELETE FROM School_Store WHERE ID = ?", itemID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	// Delete image file
	err = os.Remove("image_store/school_store/" + strconv.Itoa(itemID) + ".jpg")
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return success message
	response := restTypes.ErrorResponse{
		Code:    200,
		Message: "Item deleted from the School Store",
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
