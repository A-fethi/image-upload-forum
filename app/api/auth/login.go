package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func Login(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var potentialUser models.UserCredentials

	if err := json.NewDecoder(req.Body).Decode(&potentialUser); err != nil {
		config.Logger.Println("Login failed: Invalid request format:", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Request Format")
		return
	}

	if (potentialUser.Username == "" && potentialUser.Email == "") || potentialUser.Password == "" {
		config.Logger.Println("Login failed: Missing credentials")
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Username/Email and password are required")
		return
	}

	var (
		userID     int
		dbPassword string
		username   string
	)

	err := db.QueryRow(`
		SELECT id, password, username 
		FROM users 
		WHERE email = ? OR username = ?`,
		potentialUser.Email, potentialUser.Username,
	).Scan(&userID, &dbPassword, &username)
	if err != nil {
		if err == sql.ErrNoRows {
			config.Logger.Println("Login failed: User not found:", potentialUser.Username, potentialUser.Email)
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Error: Invalid credentials")
			return
		}
		config.Logger.Println("Login failed: Database error:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal server error")
		return
	}

	if !models.VerifyPassword(dbPassword, potentialUser.Password) {
		config.Logger.Println("Login failed: Incorrect password for user:", username)
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Error: Invalid credentials")
		return
	}

	sessionToken, err := utils.ManageSession(db, userID, username)
	if err != nil {
		config.Logger.Println("Login failed: Error generating session token for user:", username, err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal server error")
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(2 * time.Hour),
	})
	utils.SetCookie(resp, "IsLoggedIn", "true")
	config.Logger.Println("Login successful: User:", username, " User ID:", userID)

	resp.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(resp).Encode(map[string]string{"message": "Login successful"})

}
