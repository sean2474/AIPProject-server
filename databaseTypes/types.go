package databaseTypes

import "time"

// User represents a user account.
type User struct {
	ID        int    `json:"id" example:"1"`
	UserType  int    `json:"user_type" example:"2"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"johndoe@example.com"`
	Password  string `json:"-"` // exclude from Swagger docs
	RfidToken string `json:"rfid_token,omitempty" example:"RFID_TOKEN_12345"`
}

// DailySchedule represents the daily schedule of activities.
type DailySchedule struct {
	ID        int    `json:"id" example:"1"`
	Date      string `json:"date" example:"2022-01-01"`
	ImageFile []byte `json:"image_file,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	FileName  string `json:"file_name" example:"schedule.pdf"`
}

// FoodMenu represents the daily food menu.
type FoodMenu struct {
	ID        int    `json:"id" example:"1"`
	Date      string `json:"date" example:"2022-01-01"`
	Breakfast string `json:"breakfast" example:"Omelette"`
	Lunch     string `json:"lunch" example:"Pasta"`
	Dinner    string `json:"dinner" example:"Grilled chicken"`
}

// LostAndFound represents a lost and found item.
type LostAndFound struct {
	ID            int       `json:"id" example:"1"`
	ItemName      string    `json:"item_name" example:"Backpack"`
	Description   string    `json:"description,omitempty"`
	DateFound     time.Time `json:"date_found" example:"2022-01-01T12:00:00Z"`
	LocationFound string    `json:"location_found" example:"Library"`
	Status        int       `json:"status" example:"1"`
	ImageFile     []byte    `json:"image_file,omitempty"`
	SubmitterID   int       `json:"submitter_id" example:"2"`
	ImageURL      string    `json:"image_url"`
}

// SportsGame represents a sports game.
type SportsGame struct {
	ID             int       `json:"id" example:"1"`
	SportName      string    `json:"sport_name" example:"Football"`
	Category       int       `json:"category" example:"1"`
	GameLocation   string    `json:"game_location" example:"Stadium"`
	OpponentSchool string    `json:"opponent_school" example:"Central High School"`
	HomeOrAway     int       `json:"home_or_away" example:"1"`
	MatchResult    string    `json:"match_result,omitempty"`
	CoachComment   string    `json:"coach_comment,omitempty"`
	GameSchedule   time.Time `json:"game_schedule" example:"2022-01-01T12:00:00Z"`
}

// SportsInfo represents information about a sports team.
type SportsInfo struct {
	ID           int    `json:"id" example:"1"`
	SportName    string `json:"sport_name" example:"Basketball"`
	Category     int    `json:"category" example:"1"`
	Season       int    `json:"season" example:"2022"`
	CoachName    string `json:"coach_name" example:"John Smith"`
	CoachContact string `json:"coach_contact" example:"john.smith@example.com"`
	Roster       string `json:"roster,omitempty"`
}

// SchoolStore represents a product in the school store.
type SchoolStore struct {
	ID          int       `json:"ID" example:"1"`
	ProductName string    `json:"Product_Name" example:"Backpack"`
	Category    int       `json:"Category" example:"2"`
	Price       float64   `json:"Price" example:"49.99"`
	Stock       int       `json:"Stock" example:"10"`
	Description string    `json:"Description,omitempty"`
	ImageFile   []byte    `json:"image_file,omitempty"`
	DateAdded   time.Time `json:"Date_Added" example:"2022-01-01T12:00:00Z"`
}

// LoginToken represents a login token.
type LoginToken struct {
	Token  string `db:"token" json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserID int    `db:"user_id" json:"user_id" example:"1"`
}

// RfidCard represents an RFID card.
type RfidCard struct {
	Token  string `db:"token" json:"token" example:"RFID_TOKEN_12345"`
	UserID int    `db:"user_id" json:"user_id" example:"1"`
}
