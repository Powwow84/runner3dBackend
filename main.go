package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type LeaderboardEntry struct {
	Name 	string 	`json:"name"`
	Score 	int 	`json:"score"`
}

var db *sql.DB

func main(){
	fmt.Println("Server is up. Check to see if the database is responding to request")
	var err error
	db, err = sql.Open("postgres", "user=<username password=<password> dbname=<dbname host=<host> sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/leaderboard", leaderboardHandler)
	http.HandleFunc("/leaderboard/new", newEntryHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	rows ,err := db.Query("Select * FROM leaderboard")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []LeaderboardEntry 
	for rows.Next() {
		var e LeaderboardEntry 
		if err := rows.Scan(&e.Name, &e.Score); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		entries = append(entries, e)
	}

	json.NewEncoder(w).Encode(entries)
}

func newEntryHandler(w http.ResponseWriter, r *http.Request) {
	var e LeaderboardEntry
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO leaderboard (name, score) VALUES ($1, $2)", e.Name, e.Score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}