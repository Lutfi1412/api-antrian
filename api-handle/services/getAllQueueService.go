package services

import (
	"database/sql"
	"log"
	"time"
)

// Queue represents a single queue record
type Queue struct {
	ID        int `json:"id"`
	NoAntrian int `json:"no_antrian"`
	Status    int `json:"status"`
	Selected  int `json:"selected"`
}

// GetAllQueue fetches all queue data for today
func GetAllQueue(db *sql.DB) []Queue {
	// Ambil tanggal hari ini dengan zona waktu Indonesia
	loc, _ := time.LoadLocation("Asia/Jakarta")
	tanggal := time.Now().In(loc).Format("2006-01-02")

	// Query untuk PostgreSQL dengan parameter $1
	query := `SELECT id, no_antrian, status, selected FROM tbl_antrian WHERE tanggal = $1 ORDER BY id ASC`

	rows, err := db.Query(query, tanggal)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	var queues []Queue
	for rows.Next() {
		var queue Queue
		if err := rows.Scan(&queue.ID, &queue.NoAntrian, &queue.Status, &queue.Selected); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		queues = append(queues, queue)
	}

	// Memeriksa error setelah iterasi baris
	if err := rows.Err(); err != nil {
		log.Fatalf("Error during row iteration: %v", err)
	}

	return queues
}
