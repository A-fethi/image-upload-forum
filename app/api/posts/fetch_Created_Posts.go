package posts

import (
	"database/sql"
	"fmt"

	"forum/app/config"
	"forum/app/models"
)

func fetchUserCreatedPosts(posts *[]models.Post, username string, db *sql.DB) error {
	query := `
		SELECT username, id, title, content, img, Categories, created_at, likes, dislikes
		FROM posts
		WHERE username = ?
		ORDER BY created_at DESC
	`
	config.Logger.Printf("Fetching posts created by user '%s'. Executing query: %s", username, query)
	rows, err := db.Query(query, username)
	if err != nil {
		return fmt.Errorf("error querying posts created by user '%s': %v", username, err)
	}
	defer rows.Close()

	if err := rowsProcess(rows, posts, db); err != nil {
		return fmt.Errorf("error processing rows for posts created by user '%s': %v", username, err)
	}

	config.Logger.Printf("Successfully fetched posts created by user '%s'", username)
	return nil
}

func fetchUserLikedPosts(posts *[]models.Post, userID int, db *sql.DB) error {
	query := `
		SELECT posts.username, 
		posts.id AS post_id, 
		posts.title, 
		posts.content,
		posts.img,
		posts.Categories,
		posts.created_at,
		posts.likes, 
		posts.dislikes
		FROM posts
		JOIN user_interactions ON posts.id = user_interactions.item_id
		WHERE user_interactions.user_id = ? 
		AND user_interactions.item_type = "post" 
		AND user_interactions.action = "like"
		ORDER BY posts.created_at DESC
	`
	config.Logger.Printf("Fetching posts liked by user '%v'. Executing query: %s", userID, query)
	rows, err := db.Query(query, userID)
	if err != nil {
		return fmt.Errorf("error querying posts liked by user '%v': %v", userID, err)
	}
	defer rows.Close()

	if err := rowsProcess(rows, posts, db); err != nil {
		return fmt.Errorf("error processing rows for posts liked by user '%v': %v", userID, err)
	}

	config.Logger.Printf("Successfully fetched posts liked by userID '%v'", userID)
	return nil
}
