package models

import (
	"database/sql"
	"net/http"
	"regexp"
	"unicode"

	"forum/app/config"

	"golang.org/x/crypto/bcrypt"
)

type UserCredentials struct {
	Username string
	Email    string
	Password string
}

func VerifyPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

func ValidUserName(name string) (bool, string) {
	if name == "" {
		return false, "username cannot be empty"
	}
	if len(name) > 36 {
		return false, "username cannot be longer than 36 characters"
	}
	for _, char := range name {
		if !(unicode.IsLetter(char) || unicode.IsDigit(char)) {
			return false, "username contains invalid characters"
		}
	}
	return true, ""
}

func ValidEmail(email string) bool {
	valid := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	return valid.MatchString(email)
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}

func (User *UserCredentials) ValidInfo(resp http.ResponseWriter, db *sql.DB) (bool, string, int) {
	valid, message := ValidUserName(User.Username)
	if !valid {
		config.Logger.Println("Failed Register Attempt:", message)
		return false, message, http.StatusBadRequest
	}

	var username string
	err := db.QueryRow("SELECT username FROM users WHERE username=?", User.Username).Scan(&username)
	if err != nil && err != sql.ErrNoRows {
		config.Logger.Println("Database error while checking username:", err)
		return false, "internal server error", http.StatusInternalServerError
	} else if err == nil {
		return false, "username already exists", http.StatusConflict
	}

	if !ValidEmail(User.Email) {
		return false, "invalid email format", http.StatusBadRequest
	}

	var email string
	err = db.QueryRow("SELECT email FROM users WHERE email=?", User.Email).Scan(&email)
	if err != nil && err != sql.ErrNoRows {
		config.Logger.Println("Database error while checking email:", err)
		return false, "internal server error", http.StatusInternalServerError
	} else if err == nil {
		return false, "email already registered", http.StatusConflict
	}

	if !ValidatePassword(User.Password) {
		return false, "password does not meet security criteria", http.StatusBadRequest
	}

	return true, "all inputs are valid", http.StatusOK
}
