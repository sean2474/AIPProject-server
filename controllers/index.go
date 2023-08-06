package controllers

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"server/authService"
	_ "server/authService"
	"server/controllers/dailySchedule"
	"server/controllers/food"
	"server/controllers/lostAndFound"
	"server/controllers/schoolStore"
	"server/controllers/sports"
	"server/databaseControllers"
	"server/databaseTypes"
	"server/restTypes"
	"strings"
	_ "strings"
)

func writeJson(w http.ResponseWriter, resp interface{}) {
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// LoginHandler handles user authentication and generates an authentication token.
//
// @Summary Authenticate user
// @Description Login to the system and receive an authentication token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body restTypes.LoginRequest true "User login information"
// @Success 200 {object} restTypes.LoginResponse
// @Failure 401 {object} restTypes.ErrorResponse
// @Router /auth/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Check request method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req restTypes.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate credentials
	user, e := databaseControllers.GetUserByEmail(req.Username)
	if e.Code != 0 {
		writeJson(w, e)
		return
	}
	s, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	fmt.Println(string(s))

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := databaseControllers.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Build response
	resp := restTypes.LoginResponse{
		Status:  "success",
		Message: "Login successful",
		Token:   token,
		UserData: &databaseTypes.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			UserType:  user.UserType,
		},
	}
	writeJson(w, resp)

}

// TestToken greets the user with "Hello, {userName}!" if he's authorized
// @Summary Greet the user if he's authorized
// @Description Greets the user with "Hello, {userName}!" if he's authorized
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {string} string "Hello, {userName}!"
// @Failure 401 {object} restTypes.ErrorResponse "Unauthorized"
// @Failure 500 {object} restTypes.ErrorResponse "Internal Server Error"
// @Router /auth/testToken [get]
// @Router /auth/testToken [post]
func TestToken(w http.ResponseWriter, r *http.Request) {
	user, err := authService.IsAuthorized(w, r)
	if err.Code == 0 {
		fmt.Fprintf(w, "HELLO, "+user.FirstName)
		return
	}
	writeJson(w, err)

}

func FoodMenuByHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/food-menu/")
	switch r.Method {
	case "POST":
		if !authService.IsAuth(w, r) {
			return
		}
		food.PostFoodMenuHandler(w, r)
		break
	case "PUT":
		if !authService.IsAuth(w, r) {
			return
		}
		food.PutFoodMenuHandler(w, r, path)
		break
	case "DELETE":
		if !authService.IsAuth(w, r) {
			return
		}
		food.DeleteFoodMenu(w, r, path)
		break
	case "GET":
		dateStr := strings.TrimPrefix(r.URL.Path, "/data/food-menu/")
		fmt.Println(dateStr)
		if dateStr == "all" {
			food.GetAllFoodMenus(w, r)
			break
		}
		if dateStr != "" {
			food.GetFoodMenuByDate(w, r)
			break
		}
		food.GetFoodMenu(w, r)
		break
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// ScheduleImageHandler handles the image requests for the daily schedule.
func ScheduleImageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	switch r.Method {
	case "GET":
		dailySchedule.GetDailyImage(w, r)
	case "POST":
		if !authService.IsAuth(w, r) {
			return
		}
		dailySchedule.PostDailyImage(w, r)
	case "DELETE":
		if !authService.IsAuth(w, r) {
			return
		}
		dailySchedule.DeleteDailyImage(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func ScheduleHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if !authService.IsAuth(w, r) {
			return
		}
		dailySchedule.PostDailySchedule(w, r)
		break
	case "PUT":
		if !authService.IsAuth(w, r) {
			return
		}
		dailySchedule.PutDailySchedule(w, r)
		break
	case "DELETE":
		if !authService.IsAuth(w, r) {
			return
		}
		dailySchedule.DeleteDailySchedule(w, r)
		break
	case "GET":
		dailySchedule.GetDailySchedule(w, r)
		break
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func LostAndFoundHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		if !authService.IsAuth(w, r) {
			return
		}
		lostAndFound.PostLostAndFoundItem(w, r)
		break
	case "GET":
		imageID := strings.TrimPrefix(r.URL.Path, "/data/lost-and-found/image/")
		fmt.Println(imageID)
		if imageID != "/data/lost-and-found/" {
			lostAndFound.GetLostAndFoundImageHandler(w, r)
			return
		}
		lostAndFound.GetLostAndFoundItemsHandler(w, r)
	case "PUT":
		if !authService.IsAuth(w, r) {
			return
		}
		lostAndFound.PutLostAndFoundItem(w, r)
		break
	case "DELETE":
		if !authService.IsAuth(w, r) {
			return
		}
		lostAndFound.HandleDeleteLostAndFound(w, r)
	default:
		lostAndFound.GetLostAndFoundItemsHandler(w, r)
		break
	}
}

func SportsHandler(w http.ResponseWriter, r *http.Request) {
	sports.GetSportsData(w, r)
}

func GamesHandler(w http.ResponseWriter, r *http.Request) {
	sports.GetSportsGameData(w, r)
}

func SchoolStoreHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		imageID := strings.TrimPrefix(r.URL.Path, "/data/school-store/")
		fmt.Println(imageID)
		if imageID != "" {
			schoolStore.HandleSchoolStoreImage(w, r)
			return
		}
		schoolStore.HandleSchoolStore(w, r)
	case "POST":
		if !authService.IsAuth(w, r) {
			return
		}
		schoolStore.HandleAddSchoolStoreItem(w, r)
	case "PUT":
		if !authService.IsAuth(w, r) {
			return
		}
		schoolStore.PutSchoolStoreItem(w, r)
	case "DELETE":
		if !authService.IsAuth(w, r) {
			return
		}
		schoolStore.HandleDeleteSchoolStoreItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
