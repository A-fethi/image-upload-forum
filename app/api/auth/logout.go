package auth

import (
	"database/sql"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"net/http"
	"time"
)

func Logout(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	sessionToken, err := utils.GetSessionToken(req)
	if err != nil {
		if err == http.ErrNoCookie {
			config.Logger.Println("Logout failed: No session token found. User not logged in.")
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
			return
		}
		config.Logger.Println("Logout failed: Error retrieving session token:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error, Try Later.")
		return
	}

	_, username, err := utils.GetUsernameByToken(sessionToken, db)
	if err != nil {
		if err == sql.ErrNoRows {
			config.Logger.Println("Logout failed: No matching user found for the session token.")
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
			return
		}
		config.Logger.Println("Logout failed: Error retrieving username by token:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error, Try Later.")
		return
	}

	result, err := db.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	if err != nil {
		config.Logger.Println("Logout failed: Error deleting session from database:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error, Try Later.")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		config.Logger.Println("Logout failed: No active session found to delete.")
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
		return
	}

	http.SetCookie(resp, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
	})

	config.Logger.Println("Logout successful: User", username, "logged out, session removed.")
	resp.WriteHeader(http.StatusOK)
}
