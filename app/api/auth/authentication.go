package auth

import (
	"database/sql"
	"forum/app/config"
	"forum/app/models"
	"net/http"
	"strings"
)

// Authentication handles user authentication routes: register, login, logout, session check.
func Authentication(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	action := strings.TrimPrefix(req.URL.Path, "/api/auth/")

	switch action {
	case "register":
		handleAuthAction(resp, req, db, http.MethodPost, "register", Register)

	case "login":
		handleAuthAction(resp, req, db, http.MethodPost, "login", Login)

	case "logout":
		handleAuthAction(resp, req, db, http.MethodPost, "logout", Logout)

	case "session":
		if req.Method != http.MethodGet {
			config.Logger.Println("Session check failed: Invalid method:", req.Method)
			models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method Not Allowed")
			// http.Error(resp, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		if !SessionCheck(resp, req, db) {
			config.Logger.Println("Session check failed: Unauthorized access attempt")
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
			return
		}
		resp.WriteHeader(http.StatusOK)

	default:
		config.Logger.Println("Invalid authentication route attempted:", action)
		models.SendErrorResponse(resp, http.StatusNotFound, "Error: Page Not Found")
		// http.Error(resp, "404 - Page Not Found", http.StatusNotFound)
	}
}

func handleAuthAction(resp http.ResponseWriter, req *http.Request, db *sql.DB, method, action string, handler func(http.ResponseWriter, *http.Request, *sql.DB)) {
	if req.Method != method {
		config.Logger.Println("Failed attempt to", action, "using an invalid method:", req.Method)
		models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method Not Allowed")
		// http.Error(resp, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	handler(resp, req, db)
}
