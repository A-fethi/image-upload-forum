package auth

import (
	"database/sql"
	"encoding/json"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

func Register(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var (
		potentialUser models.UserCredentials
		userID        int
	)

	err := json.NewDecoder(req.Body).Decode(&potentialUser)
	if err != nil {
		config.Logger.Println("Register: Invalid request body:", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Request Format")
		return
	}

	valid, message, status := potentialUser.ValidInfo(resp, db)
	if !valid {
		models.SendErrorResponse(resp, status, message)
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", potentialUser.Email).Scan(&exists)
	if err != nil {
		config.Logger.Println("Register: Database error checking existing user:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}
	if exists {
		config.Logger.Println("Register: Email already registered:", potentialUser.Email)
		models.SendErrorResponse(resp, http.StatusConflict, "Error: Email already registered")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(potentialUser.Password), bcryptCost)
	if err != nil {
		config.Logger.Println("Register: Error hashing password:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	err = db.QueryRow(`
		INSERT INTO users (username, email, password) 
		VALUES (?, ?, ?)
		RETURNING id`,
		potentialUser.Username, potentialUser.Email, hashedPassword).Scan(&userID)
	if err != nil {
		config.Logger.Println("Register: Error inserting user into DB:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	sessionToken, err := utils.ManageSession(db, int(userID), potentialUser.Username)
	if err != nil {
		config.Logger.Println("Register: Error managing session:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(2 * time.Hour),
	})

	config.Logger.Println("Register: Successful registration for user:", potentialUser.Username, "User ID:", userID)

	resp.WriteHeader(http.StatusCreated)

}
