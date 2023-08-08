package restTypes

import (
	"server/databaseTypes"
	"time"
)

type GetEventsResponse struct {
	Events []Event `json:"events"`
}

// Event represents the data structure for the Event table
type Event struct {
	ID          string    `json:"id" example:"1"`                           // Unique identifier for the schedule
	Title       string    `json:"title" example:"New event"`                // Title of the event
	Description string    `json:"description" example:"asdasd"`             // Description of the event
	Start       time.Time `json:"start" example:"2023-08-07T04:30:00.000Z"` // Start date and time of the event (in ISO 8601 format)
	End         time.Time `json:"end" example:"2023-08-07T07:00:00.000Z"`   // End date and time of the event (in ISO 8601 format)
	Status      string    `json:"status" example:"busy"`                    // Status of the event (e.g., "busy" or "free")
	Color       string    `json:"color" example:"rgba(220,114,114,0.6)"`    // Color representation of the event
	Location    string    `json:"location" example:"sadsadad"`              // Location of the event
}

// LoginRequest represents the request body for the login API.
type LoginRequest struct {
	// User's email or username.
	//
	//
	// Required: true
	Username string `json:"username" example:"johnsmith@example.com"`

	// User's password.
	//
	// Example: mypassword123
	//
	// Required: true
	Password string `json:"password" example:"password1"`
}

// LoginResponse represents the response object returned by the login API.
type LoginResponse struct {
	// Status of the login attempt.
	//
	// Example: success
	//
	// Required: true
	Status string `json:"status" example:"success"`

	// Message indicating the result of the login attempt.
	//
	// Example: Login successful
	//
	// Required: true
	Message string `json:"message" example:"Login successful"`

	// JWT token to be used for authentication in future requests.
	//
	// Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	//
	// Required: true
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`

	// User data associated with the logged-in user.
	//
	// Required: false
	// @name User
	// @in body
	// @description User data associated with the logged-in user.
	//
	// Example: {"id":123,"first_name":"John","last_name":"Doe","email":"user@example.com","user_type":"student"}
	UserData *databaseTypes.User `json:"user_data,omitempty"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	// HTTP status code of the error response.
	//
	// Example: 400
	//
	// Required: true
	Code int `json:"code"`

	// Error message.
	//
	// Example: Invalid request
	//
	// Required: true
	Message string `json:"message"`
}

type DeleteResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AllMenuResponse struct {
	Items []databaseTypes.FoodMenu `json:"items"`
}

type LostAndFoundResponse struct {
	Items []databaseTypes.LostAndFound `json:"items"`
}

type LostAndFoundPostResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      int64  `json:"id,omitempty"`
}

type LostAndFoundErrorResponse struct {
	Error string `json:"error"`
}

type LostAndFoundInput struct {
	ItemName      string `json:"item_name"`
	Description   string `json:"description"`
	DateFound     string `json:"date_found"`
	LocationFound string `json:"location_found"`
	Status        string `json:"status"`
}

type SportsDataList struct {
	List []databaseTypes.SportsInfo `json:"list"`
}

type SportsGameDataList struct {
	List []databaseTypes.SportsGame `json:"list"`
}

type SchoolStoreResponse struct {
	List []databaseTypes.SchoolStore `json:"list"`
}

type SchoolStoreErrorResponse struct {
	Error string `json:"error"`
}

type SchoolStorePostResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}
