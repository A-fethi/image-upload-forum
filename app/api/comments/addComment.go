package comments

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func AddComment(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var CommentID int
	var Comment struct {
		Content string `json:"content"`
		Post_id int    `json:"post_id"`
	}

	err := json.NewDecoder(req.Body).Decode(&Comment)
	if err != nil {
		config.Logger.Printf("Error decoding request body: %v\n", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Content")

		return
	}

	session_token, err := utils.GetSessionToken(req)
	if err != nil || session_token == "" || !auth.SessionCheck(resp, req, db) {
		config.Logger.Println("User not authenticated: ", err)
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")

		return
	}

	if Comment.Content == "" || Comment.Post_id == 0 {
		config.Logger.Println("Comment content, username cannot be empty")
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Content")

		return
	}

	var username string
	var userID int
	userID, username, err = utils.GetUsernameByToken(session_token, db)
	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")

		return
	}
	config.Logger.Println("Username: ", username, "UserID: ", userID, "Comment: ", Comment)
	err = db.QueryRow(`
		INSERT INTO comments (post_id, user_id, author, content, created_at)
		VALUES (?, ?, ?, ?, ?)
		RETURNING id`,
		Comment.Post_id, userID, username, Comment.Content, time.Now().Format("2006-01-02 15:04:05")).Scan(&CommentID)
	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")

		return
	}

	config.Logger.Println("Comment add: ", Comment)
	comm := models.Comment{
		ID:        CommentID,
		Content:   Comment.Content,
		Username:  username,
		PostID:    Comment.Post_id,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		Likes:     0,
		Dislikes:  0,
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
	json.NewEncoder(resp).Encode(comm)
}
