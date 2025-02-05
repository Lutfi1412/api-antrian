package services

import (
	"database/sql"
	"time"
)

// GetJumlahAntrian mengambil jumlah antrian untuk tanggal hari ini
func GetJumlahAntrian(db *sql.DB) int {
	// Ambil tanggal hari ini dengan zona waktu Indonesia
	loc, _ := time.LoadLocation("Asia/Jakarta")
	tanggal := time.Now().In(loc).Format("2006-01-02")

	// Query untuk menghitung jumlah antrian berdasarkan tanggal
	var jumlahAntrian int
	query := `SELECT count(id) FROM tbl_antrian WHERE tanggal = $1`
	db.QueryRow(query, tanggal).Scan(&jumlahAntrian)
	return jumlahAntrian
}
