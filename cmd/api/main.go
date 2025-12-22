package main

import (
	"log"
	"net/http"
	"os"

	"thainsbook/internal/handlers"
	"thainsbook/internal/models"
	"thainsbook/internal/utils"

	"github.com/joho/godotenv"
)

func main() {
	// Load Env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	prefix := os.Getenv("PREFIX")
	jwt := os.Getenv("JWT_KEY")

	// Load DB Connection
	db, ok := utils.ConnectDB()
	if ok != nil {
		log.Fatalf("Error connecting to database: %s", ok)
	}
	defer db.Close()

	// Load Dependency Injection
	app := handlers.Application{
		Users:   models.UserModel{DB: db},
		Entries: models.EntryModel{DB: db},
		JWT:     jwt,
	}

	// Initialize Server
	r := http.NewServeMux()

	// E.g. endpoint: /api/v1/entries
	// User Endpoints
	r.HandleFunc("POST "+prefix+"/users/register", app.HandleRegister)
	r.HandleFunc("POST "+prefix+"/users/login", app.HandleLogin)

	// Entries Endpoints, specific to user, need to authenticate
	r.HandleFunc("GET "+prefix+"/entries", app.Authenticate(app.HandleGetUserEntries))
	r.HandleFunc("GET "+prefix+"/entries/{id}", app.Authenticate(app.HandleGetUserEntry))
	r.HandleFunc("POST "+prefix+"/entries", app.Authenticate(app.HandleCreateEntry))
	r.HandleFunc("PATCH "+prefix+"/entries/{id}", app.Authenticate(app.HandleUpdateEntry))
	r.HandleFunc("DELETE "+prefix+"/entries/{id}", app.Authenticate(app.HandleDeleteEntry))

	// Catchall 404
	r.HandleFunc("/", handlers.HandleNotFound)

	// Start Server
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.Logger(r)))
}
