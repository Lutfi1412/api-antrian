package services

import (
	"database/sql"
	"log"
)

// DeleteDate menghapus data dari tabel `tbl_antrian` berdasarkan tanggal
func DeleteDate(db *sql.DB, tanggal string) {
	// Query untuk PostgreSQL dengan parameter $1
	query := "DELETE FROM tbl_antrian WHERE tanggal = $1"
	result, err := db.Exec(query, tanggal)
	if err != nil {
		log.Fatalf("Error executing delete query: %v", err)
	}

	// Mengecek jumlah baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Error getting rows affected: %v", err)
	}

	// Jika tidak ada baris yang terhapus, log informasi
	if rowsAffected == 0 {
		log.Println("No rows were deleted")
	} else {
		log.Printf("%d rows deleted successfully", rowsAffected)
	}
}
