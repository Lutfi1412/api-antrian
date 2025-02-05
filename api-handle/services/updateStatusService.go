package services

import (
	"database/sql"
	"log"
	"time"
)

// UpdateAntrianService memperbarui status antrian berdasarkan ID
func UpdateAntrianService(db *sql.DB, id int, status string, loket string, selected string) {
	tanggal := time.Now().Add(7 * time.Hour).Format("2006-01-02 15:04:05")

	// Debug log: data yang akan digunakan
	log.Printf("UpdateAntrianService: id=%d, status=%s, loket=%s, selected=%s, tanggal=%s\n", id, status, loket, selected, tanggal)

	// Eksekusi query
	result, err := db.Exec(`UPDATE tbl_antrian SET status=$1, updated_date=$2, loket=$3, selected=$4 WHERE id=$5`,
		status, tanggal, loket, selected, id)
	if err != nil {
		log.Printf("Error executing UPDATE query: %v", err)
		return
	}

	// Debug log: jumlah baris yang terpengaruh
	rowsAffected, _ := result.RowsAffected()
	log.Printf("UpdateAntrianService: rows affected=%d", rowsAffected)
}
