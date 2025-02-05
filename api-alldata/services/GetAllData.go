package services

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func GetAllData(db *sql.DB) (map[string]interface{}, error) {
	// Set Jakarta timezone
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	tanggal := time.Now().In(loc).Format("2006-01-02")

	// Database query to fetch data
	query := `SELECT
		(SELECT count(id) FROM tbl_antrian WHERE tanggal = $1) AS jumlah_antrian,
		(SELECT no_antrian FROM tbl_antrian WHERE tanggal = $1 AND status = '0' ORDER BY no_antrian ASC LIMIT 1) AS antrian_selanjutnya,
		(SELECT count(id) FROM tbl_antrian WHERE tanggal = $1 AND status = '0') AS sisa_antrian,
		(SELECT no_antrian FROM tbl_antrian WHERE tanggal = $1 AND status = '2' AND loket = '1' ORDER BY updated_date DESC LIMIT 1) AS antrian_loket_1,
		(SELECT no_antrian FROM tbl_antrian WHERE tanggal = $1 AND status = '2' AND loket = '2' ORDER BY updated_date DESC LIMIT 1) AS antrian_loket_2,
		(SELECT no_antrian FROM tbl_antrian WHERE tanggal = $1 AND status = '2' AND loket = '3' ORDER BY updated_date DESC LIMIT 1) AS antrian_loket_3`

	var jumlahAntrian int
	var antrianSelanjutnya sql.NullString
	var sisaAntrian int
	var antrianLoket1 sql.NullString
	var antrianLoket2 sql.NullString
	var antrianLoket3 sql.NullString

	// Query the database
	err := db.QueryRow(query, tanggal).Scan(
		&jumlahAntrian,
		&antrianSelanjutnya,
		&sisaAntrian,
		&antrianLoket1,
		&antrianLoket2,
		&antrianLoket3,
	)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	// Prepare data for response
	data := map[string]interface{}{
		"jumlah_antrian":      jumlahAntrian,
		"antrian_selanjutnya": nullStringToString(antrianSelanjutnya),
		"sisa_antrian":        sisaAntrian,
		"antrian_loket_1":     nullStringToString(antrianLoket1),
		"antrian_loket_2":     nullStringToString(antrianLoket2),
		"antrian_loket_3":     nullStringToString(antrianLoket3),
	}

	return data, nil
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "0"
}
