package main

import (
	"github.com/rs/cors" // Import the cors package
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"server/controllers"
	_ "server/docs"
)

// Rest of your code...
// @title Go Rest API with Swagger for school system
// @version 1.0
// @description Simple swagger implementation in Go HTTP
// @contact.name Senya
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Create a new cors handler with permissive options (allowing all origins)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins, you can restrict this to specific origins if needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Apply the cors handler to your existing handlers
	http.Handle("/swagger/", corsHandler.Handler(httpSwagger.WrapHandler))
	http.Handle("/auth/login", corsHandler.Handler(http.HandlerFunc(controllers.LoginHandler)))
	//http.Handle("/auth/testToken", corsHandler.Handler(http.HandlerFunc(controllers.TestToken)))
	//http.Handle("/auth/testToken", corsHandler.Handler(http.HandlerFunc(controllers.SchoolStoreHandler)))
	http.Handle("/data/food-menu/", corsHandler.Handler(http.HandlerFunc(controllers.FoodMenuByHandler)))
	http.Handle("/data/daily-schedule/image", corsHandler.Handler(http.HandlerFunc(controllers.ScheduleImageHandler)))
	http.Handle("/data/daily-schedule/", corsHandler.Handler(http.HandlerFunc(controllers.ScheduleHandler)))
	http.Handle("/data/lost-and-found/", corsHandler.Handler(http.HandlerFunc(controllers.LostAndFoundHandler)))
	http.Handle("/data/sports/", corsHandler.Handler(http.HandlerFunc(controllers.SportsHandler)))
	http.Handle("/data/games/", corsHandler.Handler(http.HandlerFunc(controllers.GamesHandler)))
	http.Handle("/data/school-store/", corsHandler.Handler(http.HandlerFunc(controllers.SchoolStoreHandler)))

	// Start the server with your handlers
	http.ListenAndServe(":8082", nil)
}
