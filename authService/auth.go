package authService

import (
	"database/sql"
	"fmt"
	"net/http"
	"server/databaseTypes"
	"server/restTypes"
	"strings"
	"time"
)

func IsAuth(w http.ResponseWriter, r *http.Request) bool {
	// Check authentication
	_, erro := IsAuthorized(w, r)
	if erro.Code != 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func IsAuthorized(w http.ResponseWriter, r *http.Request) (databaseTypes.User, restTypes.ErrorResponse) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Get the Authorization header from the request
	authHeader := r.Header.Get("Authorization")
	// Check if the Authorization header is present and has the correct format
	if authHeader == "" {
		// If the Authorization header is missing, return error
		return databaseTypes.User{}, restTypes.ErrorResponse{
			Message: "Authorization header is missing",
			Code:    401,
		}
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		// If the Authorization header has an invalid format, return error
		return databaseTypes.User{}, restTypes.ErrorResponse{
			Message: "Authorization header has an invalid format",
			Code:    401,
		}
	}

	// Extract the token from the Authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Query the database for the most recent token for the given token
	var addedAt time.Time
	var user databaseTypes.User
	err = db.QueryRow("SELECT id, user_type, first_name,last_name, email, added_at FROM LoginTokens, Users WHERE (LoginTokens.user_id = Users.id AND token = ? )ORDER BY added_at DESC LIMIT 1", token).Scan(&user.ID, &user.UserType, &user.FirstName, &user.LastName, &user.Email, &addedAt)
	fmt.Println(err)
	if err != nil {
		if err == sql.ErrNoRows {
			time.Sleep(2 * time.Second)
			// If the token is not found in the database, return error
			return databaseTypes.User{}, restTypes.ErrorResponse{
				Message: "Token is not found in the database",
				Code:    401,
			}
		}
		// If there was an error querying the database, return error
		return databaseTypes.User{}, restTypes.ErrorResponse{
			Message: "Internal Server Error",
			Code:    500,
		}
	}

	// Check if the token is older than one day
	if time.Since(addedAt) > 24*time.Hour {
		// If the token is older than one day, delete it from the database
		_, err = db.Exec("DELETE FROM LoginTokens WHERE token = ?", token)
		if err != nil {
			// If there was an error deleting the token from the database, return error
			return databaseTypes.User{}, restTypes.ErrorResponse{
				Message: "Internal Server Error",
				Code:    500,
			}
		}
		// Return false to indicate that the token was deleted
		return databaseTypes.User{}, restTypes.ErrorResponse{
			Message: "Token is too old",
			Code:    401,
		}
	}
	return user, restTypes.ErrorResponse{Code: 0}
}
