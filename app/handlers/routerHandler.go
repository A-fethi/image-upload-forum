package handlers

import (
	"database/sql"
	"forum/app/api/auth"
	"forum/app/api/comments"
	"forum/app/api/posts"
	"forum/app/api/reactions"
	"forum/app/models"
	"net/http"
	"strings"
)

func Router(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	path := strings.Split(req.URL.Path[5:], "/")
	if len(path) == 0 || path[0] == "" {
		models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
		return
	}

	switch strings.ToLower(path[0]) {
	case "auth":
		auth.Authentication(resp, req, db)
		return

	case "posts":
		if req.Method == http.MethodPost && len(path) > 1 && path[1] == "add" {
			posts.AddPost(resp, req, db)
			return
		}
		if req.Method == http.MethodGet {
			posts.GetPosts(resp, req, db)
			return
		}
		models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method not allowed")
		return

	case "comments":
		if req.Method == http.MethodPost {
			comments.AddComment(resp, req, db)
			return
		}
		models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Error: Method not allowed")
		return

	case "reactions":
		reactions.AddReaction(resp, req, db)
		return

	default:
		models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
		return
	}
}
