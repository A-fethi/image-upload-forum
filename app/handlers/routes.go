package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"forum/app/utils"
)

func RegisterRoutes(DB *sql.DB) {
	http.HandleFunc("/static/", Static)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, DB)
	})
	
	http.HandleFunc("/api/", utils.RateLimitMiddleware(func(w http.ResponseWriter, r *http.Request) {
		Router(w, r, DB)
	}, 5, 1*time.Second))
}
