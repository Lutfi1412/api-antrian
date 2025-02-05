package services

import (
	"database/sql"
	"log"
)

// DateWithLoketCount struct menyimpan data tanggal dan jumlah data per loket
type DateWithLoketCount struct {
	Tanggal string `json:"tanggal"`
	Jumlah  int    `json:"jumlah"`
	Loket1  int    `json:"loket1"`
	Loket2  int    `json:"loket2"`
	Loket3  int    `json:"loket3"`
}

// GetDateService mengambil tanggal unik dan jumlah data per loket
func GetDateService(db *sql.DB) []DateWithLoketCount {
	// Query untuk PostgreSQL
	query := `
        SELECT 
            tanggal,
            COUNT(id) AS jumlah,
            SUM(CASE WHEN loket = '1' THEN 1 ELSE 0 END) AS loket1,
            SUM(CASE WHEN loket = '2' THEN 1 ELSE 0 END) AS loket2,
            SUM(CASE WHEN loket = '3' THEN 1 ELSE 0 END) AS loket3
        FROM tbl_antrian
        GROUP BY tanggal
        ORDER BY tanggal;
    `

	// Menjalankan query
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	// Menyiapkan hasil untuk dikembalikan
	var dates []DateWithLoketCount
	for rows.Next() {
		var date DateWithLoketCount
		if err := rows.Scan(&date.Tanggal, &date.Jumlah, &date.Loket1, &date.Loket2, &date.Loket3); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		dates = append(dates, date)
	}

	// Memeriksa jika ada error pada hasil query
	if err := rows.Err(); err != nil {
		log.Fatalf("Row iteration error: %v", err)
	}

	return dates
}
