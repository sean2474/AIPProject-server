package databaseControllers

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"server/databaseTypes"
	"server/restTypes"
	"time"
)

func GetTokenForUser(userID int) (string, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// Query the database for the most recent token for the given user ID
	var token string
	var addedAt time.Time
	err = db.QueryRow("SELECT token, added_at FROM LoginTokens WHERE user_id = ? ORDER BY added_at DESC LIMIT 1", userID).Scan(&token, &addedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// If there are no tokens for the given user ID, return an empty string and nil error
			return "", nil
		}
		// If there was an error querying the database, return an error
		return "", err
	}

	// Check if the token is older than one day
	if time.Since(addedAt) > 24*time.Hour {
		// If the token is older than one day, delete it from the database
		_, err = db.Exec("DELETE FROM LoginTokens WHERE token = ?", token)
		if err != nil {
			return "", err
		}
		// Return an empty string and nil error to indicate that the token was deleted
		return "", nil
	}

	// If the token is not older than one day, return it
	return token, nil
}

func GetUserByEmail(email string) (*databaseTypes.User, restTypes.ErrorResponse) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var user databaseTypes.User
	row := db.QueryRow("SELECT * FROM Users WHERE email = ?", email)
	err = row.Scan(&user.ID, &user.UserType, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, restTypes.ErrorResponse{Message: "user not found", Code: 401}
		}
		fmt.Println(err.Error())
		return nil, restTypes.ErrorResponse{Message: "FATAL", Code: 500}
	}
	return &user, restTypes.ErrorResponse{Code: 0}
}

func GenerateToken(userID int) (string, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	t, _ := GetTokenForUser(userID)
	if t != "" {
		return t, nil
	}
	token := uuid.New().String()

	// Insert the token and user_id into the LoginTokens table
	_, err = db.Exec("INSERT INTO LoginTokens (token, user_id) VALUES (?, ?)", token, userID)
	if err != nil {
		return "", fmt.Errorf("error inserting token into database: %w", err)
	}

	// Return the generated token
	return token, nil
}
