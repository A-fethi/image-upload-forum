package utils

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"forum/app/config"

	"github.com/gofrs/uuid"
)

func GetSessionToken(req *http.Request) (string, error) {
	cookie, err := req.Cookie("session_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func ManageSession(db *sql.DB, userID int, username string) (string, error) {
	token, err := uuid.NewV7()
	if err != nil {
		config.Logger.Println("Error generating session token: ", err)
		return "", err
	}

	var sessionExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE user_id = ?)", userID).Scan(&sessionExists)
	if err != nil {
		config.Logger.Println("Database error checking session existence: ", err)
		return "", err
	}

	if sessionExists {
		_, err = db.Exec(
			"UPDATE sessions SET session_token = ?, isloggedin = ?, created_at = ? WHERE user_id = ?",
			token.String(), true, time.Now(), userID,
		)
		if err != nil {
			config.Logger.Println("Failed to update session for user ID:", userID, " Error:", err)
			return "", err
		}
	} else {
		_, err = db.Exec(
			"INSERT INTO sessions (user_id, username, session_token, isloggedin) VALUES (?, ?, ?, ?)",
			userID, username, token.String(), true,
		)
		if err != nil {
			config.Logger.Println("Failed to create session for user ID:", userID, " Error:", err)
			return "", err
		}
	}

	return token.String(), nil
}

func LoggedInUser(db *sql.DB, token string) bool {
	var isLoggedIn bool
	err := db.QueryRow("SELECT isloggedin FROM sessions WHERE session_token = ?", token).Scan(&isLoggedIn)
	if err != nil {
		if err == sql.ErrNoRows {
			config.Logger.Println("Invalid token")
			return false
		}
		config.Logger.Println("Internal Server Error: ", err)
		return false
	}

	return isLoggedIn
}

func GetUsernameByToken(token string, db *sql.DB) (int, string, error) {

	var username string
	var userID int
	err := db.QueryRow("SELECT user_id, username FROM sessions WHERE session_token = ?", token).Scan(&userID, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", errors.New("invalid token")
		}
		return 0, "", errors.New("internal server error")
	}

	return userID, username, nil
}

func SetCookie(resp http.ResponseWriter, name, value string) {
	config.Logger.Println("Setting cookie:", name, "=", value)
	http.SetCookie(resp, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	})
}
