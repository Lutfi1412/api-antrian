package handler

import (
	"database/sql"
	"fmt"
	"log"
	"myapp/services"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	. "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Buat server G4Vercel
	server := New()
	server.Use(Recovery(func(err interface{}, c *Context) {
		if httpError, ok := err.(HttpError); ok {
			c.JSON(httpError.Status, H{
				"message": httpError.Error(),
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, H{
				"message": message,
			})
		}
	}))
	server.GET("/api/alldata", GetAllDataHandler)

	// Tambahkan CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                       // Domain React Anda
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Metode yang diizinkan
		AllowedHeaders:   []string{"Content-Type", "Authorization"},           // Header yang diizinkan
		AllowCredentials: true,                                                // Izinkan kredensial seperti cookies
	})

	// Bungkus handler utama dengan CORS
	corsWrappedHandler := corsHandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handle(w, r)
	}))

	// Panggil handler yang sudah dibungkus
	corsWrappedHandler.ServeHTTP(w, r)
}

func GetAllDataHandler(c *Context) {

	// Initialize database connection
	db := InitializeDB()
	defer db.Close()

	// Get data from the database service
	data, err := services.GetAllData(db)
	if err != nil {
		// If there is an error, send a 500 error with the message
		c.JSON(500, H{
			"message": err.Error(),
		})
		return
	}

	// Set the response type to JSON and send the data back
	c.JSON(200, H{
		"data": data,
	})
}

func InitializeDB() *sql.DB {
	// Read database connection credentials from environment variables
	dbUser := os.Getenv("DB_USER")         // Your database user
	dbPassword := os.Getenv("DB_PASSWORD") // Your database password
	dbHost := os.Getenv("DB_HOST")         // Your database host
	dbPort := os.Getenv("DB_PORT")         // Your database port
	dbName := os.Getenv("DB_NAME")         // Your database name

	// If any of these environment variables are missing, log an error
	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("One or more environment variables are missing. Please check your environment setup.")
	}

	// Create connection string using the environment variables
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open connection to PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	// Return the database connection object
	return db
}
