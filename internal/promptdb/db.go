package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=koyeb-adm password=npg_q2XJKtRBLC4l host=ep-soft-waterfall-a44o6qqo.us-east-1.pg.koyeb.app dbname=koyebdb sslmode=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create tapsq
		log.Fatal(err)
	}

	// Insert example
	_, err = db.Exec(`INSERT INTO api_keys (key_prefix, key_hash, owner_id) VALUES ($1, $2, $3)`,
		"abc123", "hashvalue", "user42")
	if err != nil {
		log.Fatal(err)
	}

	// Query example
	rows, err := db.Query(`SELECT id, owner_id FROM api_keys`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var owner string
		if err := rows.Scan(&id, &owner); err != nil {
			log.Fatal(err)
		}
		log.Println("ID:", id, "Owner:", owner)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
