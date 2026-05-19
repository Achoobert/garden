package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	TimeBudget int    `json:"time_budget"` // minutes per day
	CashBudget float64 `json:"cash_budget"`
	SpaceType  string  `json:"space_type"`
}

type Plant struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	WaterFreq    int    `json:"water_freq"` // minutes
	FeedFreq     int    `json:"feed_freq"`
	PestFreq     int    `json:"pest_freq"`
	SunFreq      int    `json:"sun_freq"`
}

type UserPlant struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	PlantID   int       `json:"plant_id"`
	PlantName string    `json:"plant_name"`
	LastWater time.Time `json:"last_water"`
	LastFeed  time.Time `json:"last_feed"`
	LastPest  time.Time `json:"last_pest"`
	LastSun   time.Time `json:"last_sun"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./igarden.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, time_budget INTEGER, cash_budget REAL, space_type TEXT);
	CREATE TABLE IF NOT EXISTS plants (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, water_freq INTEGER, feed_freq INTEGER, pest_freq INTEGER, sun_freq INTEGER);
	CREATE TABLE IF NOT EXISTS user_plants (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, plant_id INTEGER, last_water DATETIME, last_feed DATETIME, last_pest DATETIME, last_sun DATETIME);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// Seed plants
	db.Exec("INSERT OR IGNORE INTO plants (id, name, water_freq, feed_freq, pest_freq, sun_freq) VALUES (1, 'Holy Basil (Kraphao)', 5, 10, 15, 20)")
	db.Exec("INSERT OR IGNORE INTO plants (id, name, water_freq, feed_freq, pest_freq, sun_freq) VALUES (2, 'Bird''s Eye Chili (Phrik Ki Nu)', 5, 12, 18, 25)")

}

func main() {
	initDB()
	defer db.Close()

	http.Handle("/", http.FileServer(http.Dir("./web/static")))
	
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/garden", gardenHandler)
	http.HandleFunc("/api/action", actionHandler)

	fmt.Println("Server starting on :30022")
	log.Fatal(http.ListenAndServe(":30022", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	
	res, err := db.Exec("INSERT INTO users (username, time_budget, cash_budget, space_type) VALUES (?, ?, ?, ?)", 
		u.Username, u.TimeBudget, u.CashBudget, u.SpaceType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)

	// Auto-add the two PoC plants
	db.Exec("INSERT INTO user_plants (user_id, plant_id, last_water, last_feed, last_pest, last_sun) VALUES (?, 1, ?, ?, ?, ?)", id, time.Now(), time.Now(), time.Now(), time.Now())
	db.Exec("INSERT INTO user_plants (user_id, plant_id, last_water, last_feed, last_pest, last_sun) VALUES (?, 2, ?, ?, ?, ?)", id, time.Now(), time.Now(), time.Now(), time.Now())

	json.NewEncoder(w).Encode(u)
}

func gardenHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	rows, err := db.Query(`
		SELECT up.id, up.user_id, up.plant_id, p.name, up.last_water, up.last_feed, up.last_pest, up.last_sun 
		FROM user_plants up 
		JOIN plants p ON up.plant_id = p.id 
		WHERE up.user_id = ?`, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var plants []UserPlant
	for rows.Next() {
		var p UserPlant
		rows.Scan(&p.ID, &p.UserID, &p.PlantID, &p.PlantName, &p.LastWater, &p.LastFeed, &p.LastPest, &p.LastSun)
		plants = append(plants, p)
	}
	json.NewEncoder(w).Encode(plants)
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserPlantID int    `json:"user_plant_id"`
		Action      string `json:"action"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	field := ""
	switch req.Action {
	case "water": field = "last_water"
	case "feed": field = "last_feed"
	case "pest": field = "last_pest"
	case "sun": field = "last_sun"
	}

	if field != "" {
		db.Exec(fmt.Sprintf("UPDATE user_plants SET %s = ? WHERE id = ?", field), time.Now(), req.UserPlantID)
	}
	w.WriteHeader(http.StatusOK)
}
