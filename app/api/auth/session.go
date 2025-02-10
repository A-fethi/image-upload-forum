package auth

import (
	"database/sql"
	"forum/app/config"
	"forum/app/utils"
	"net/http"
	"time"
)


func SessionCheck(resp http.ResponseWriter, req *http.Request, db *sql.DB) bool {
	session_token, err := utils.GetSessionToken(req)
	if err != nil || session_token == "" {
		config.Logger.Println("SessionCheck: No session token found")
		return false
	}

	var createdAt time.Time
	err = db.QueryRow("SELECT created_at FROM sessions WHERE session_token = ? AND isloggedin = 1", session_token).Scan(&createdAt)
	if err == sql.ErrNoRows {
		config.Logger.Println("SessionCheck: No active session found for token:", session_token)
		return false
	} else if err != nil {
		config.Logger.Println("SessionCheck: Database error:", err)
		return false
	}

	if time.Since(createdAt) > 2*time.Hour {
		config.Logger.Println("SessionCheck: Session token expired, deleting from database")
		_, delErr := db.Exec("DELETE FROM sessions WHERE session_token = ?", session_token)
		if delErr != nil {
			config.Logger.Println("SessionCheck: Failed to delete expired session:", delErr)
		}
		return false
	}

	config.Logger.Println("SessionCheck: Session is valid for token:", session_token)
	return true
}
