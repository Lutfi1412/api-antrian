package services

import (
	"database/sql"
)

// GetNextAntrian mendapatkan nomor antrian berikutnya untuk hari ini
func GetNextAntrian(db *sql.DB) int {
	var noAntrian int
	var result sql.NullInt64 // Gunakan sql.NullInt64 untuk menangani nilai NULL

	// Menggunakan CURRENT_DATE untuk mendapatkan tanggal hari ini di PostgreSQL
	query := "SELECT max(no_antrian) FROM tbl_antrian WHERE tanggal = CURRENT_DATE"
	db.QueryRow(query).Scan(&result)

	if result.Valid {
		noAntrian = int(result.Int64) + 1
	} else {
		noAntrian = 1
	}
	return noAntrian
}

// InsertAntrian menambahkan data antrian baru ke dalam tabel
func InsertAntrian(db *sql.DB, tanggal string, noAntrian int) {
	// Menggunakan placeholder $1, $2 untuk PostgreSQL
	query := `INSERT INTO tbl_antrian (tanggal, no_antrian) VALUES ($1, $2)`
	db.Exec(query, tanggal, noAntrian)
}
