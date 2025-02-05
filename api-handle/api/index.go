package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"my-app-backend/services"
	"net/http"
	"os"
	"time"

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

	// Gunakan metode yang tepat untuk operasi
	server.GET("/api/alldetail", GetAllQueueHandler)       // GET karena mengambil data
	server.POST("/api/update", UpdateStatusHandler)        // PUT karena update data
	server.POST("/api/delete", DeleteDateHandler)          // DELETE karena menghapus data
	server.GET("/api/date", GetDateQueueHandler)           // GET karena mengambil data
	server.GET("/api/jumlahantrian", JumlahAntrianHandler) // GET karena mengambil data
	server.POST("/api/insert", AntrianHandler)             // POST karena menyisipkan data
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

type QueueResponse struct {
	Queues []services.Queue `json:"queues"`
}

// Handler to get all the queue data
func GetAllQueueHandler(c *Context) {

	// Initialize the database connection
	db := InitializeDB()
	defer db.Close()

	// Panggil fungsi untuk mendapatkan data antrian
	queues := services.GetAllQueue(db)

	// Persiapkan data response
	response := QueueResponse{Queues: queues}

	c.JSON(200, H{
		"data": response,
	})
}

type DeleteDateRequest struct {
	Tanggal string `json:"tanggal"`
}

// DeleteDateHandler menangani penghapusan data berdasarkan tanggal
func DeleteDateHandler(c *Context) {
	// Dekode body request
	var reqBody DeleteDateRequest
	json.NewDecoder(c.Req.Body).Decode(&reqBody)

	// Inisialisasi koneksi database
	db := InitializeDB()
	defer db.Close()

	services.DeleteDate(db, reqBody.Tanggal)

	// Jika penghapusan berhasil, kirimkan response sukses
	c.JSON(200, H{
		"message": "Data successfully deleted",
	})
}

type DateWithLoketResponse struct {
	Dates []services.DateWithLoketCount `json:"dates"`
}

// Handler untuk mendapatkan data tanggal dan jumlah per loket
func GetDateQueueHandler(c *Context) {

	// Inisialisasi koneksi ke database
	db := InitializeDB()
	defer db.Close()

	// Panggil fungsi dari service
	dates := services.GetDateService(db)

	// Persiapkan response
	response := DateWithLoketResponse{Dates: dates}

	c.JSON(200, H{
		"data": response,
	})
}

func JumlahAntrianHandler(c *Context) {

	db := InitializeDB()
	defer db.Close()

	jumlahAntrian := services.GetJumlahAntrian(db)

	// Mengembalikan data sebagai JSON
	response := map[string]string{
		"jumlah_antrian": fmt.Sprintf("%d", jumlahAntrian),
	}

	c.JSON(200, H{
		"data": response,
	})
}

func AntrianHandler(c *Context) {
	// Koneksi ke database
	db := InitializeDB()
	defer db.Close()
	noAntrian := services.GetNextAntrian(db)
	// Ambil tanggal hari ini dengan zona waktu Asia/Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Error loading timezone: %v", err)
		c.JSON(500, H{
			"error": "Failed to load timezone",
		})
		return
	}
	tanggal := time.Now().In(loc).Format("2006-01-02")
	services.InsertAntrian(db, tanggal, noAntrian)
	response := map[string]string{"status": "Sukses"}

	c.JSON(200, H{
		"data": response,
	})
}

func UpdateStatusHandler(c *Context) {
	// Dekode body permintaan
	var requestData struct {
		ID       int    `json:"id"`
		Status   string `json:"status"`
		Loket    string `json:"loket"`
		Selected string `json:"selected"`
	}
	json.NewDecoder(c.Req.Body).Decode(&requestData)
	// Inisialisasi koneksi database
	db := InitializeDB()
	defer db.Close()

	// Panggil service untuk update data
	services.UpdateAntrianService(db, requestData.ID, requestData.Status, requestData.Loket, requestData.Selected)

	// Response berhasil
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Status antrian berhasil diperbarui",
	}

	// Kirimkan response sukses
	c.JSON(http.StatusOK, H{
		"data": response,
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
