package handlers

import (
	"database/sql"
	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"net/http"
)

func Home(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	config.Logger.Println("Handling request for Home page")

	if req.URL.Path != "/" {
		config.Logger.Println("404 error - Page not found for path:", req.URL.Path)
		models.SendErrorResponse(resp, http.StatusNotFound, "Error: Page Not Found")
		return
	}

	config.Logger.Println("Checking for session token")
	sessionToken, err := utils.GetSessionToken(req)
	if err != nil {
		config.Logger.Println("No session token found, user not logged in")
		config.Templates.ExecuteTemplate(resp, "home.html", nil)
		SetCookie(resp, "session_token", "")
		return
	}

	config.Logger.Println("Session token found, verifying...")
	var username string

	isLoggedIn := auth.SessionCheck(resp, req, db)
	err = db.QueryRow("SELECT username FROM sessions WHERE session_token = ?", sessionToken).Scan(&username)
	if err != nil || !isLoggedIn {
		config.Logger.Println("Invalid or expired session token for token:", sessionToken)
		SetCookie(resp, "IsLoggedIn", "false")
		SetCookie(resp, "Username", "")
		SetCookie(resp, "session_token", "")
		config.Templates.ExecuteTemplate(resp, "home.html", nil)
		return
	}
	User := config.TemplateData{
		IsAuthenticated: true,
		Username:        username,
	}
	config.Logger.Println("User is logged in:", User)
	config.Logger.Println("User successfully logged in:", username)
	config.Templates.ExecuteTemplate(resp, "home.html", User)
	SetCookie(resp, "session_token", sessionToken)
}

func SetCookie(resp http.ResponseWriter, name, value string) {
	config.Logger.Println("Setting cookie:", name, "=", value)
	http.SetCookie(resp, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
	})
}
