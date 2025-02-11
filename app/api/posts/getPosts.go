package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/app/api/auth"
	"forum/app/api/comments"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func GetPosts(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var posts []models.Post
	var categories []string
	var err error

	page := 1
	switch {
	case strings.HasPrefix(req.URL.Path, "/api/posts/categories="):
		pathParts := strings.Split(req.URL.Path, "/")
		catPart := strings.Split(pathParts[3], "categories=")
		if len(pathParts) == 5 {
			if p, err := strconv.Atoi(pathParts[4]); err == nil && p > 0 {
				page = p
			} else {
				models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
				return
			}

		} else {
			models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
			return
		}
		if len(catPart) == 2 {
			categories = strings.Split(catPart[1], "&")
			config.Logger.Printf("Categories found in URL: %v", categories)
		}
		if len(categories) > 0 {
			config.Logger.Printf("Fetching posts by categories: %v for page %d", categories, page)
			err = fetchPostsByCategories(categories, &posts, page, db)
		} else {
			models.SendErrorResponse(resp, http.StatusBadRequest, "Invalid categories parameter")
			return
		}

	case req.URL.Path == "/api/posts/created":
		if !auth.SessionCheck(resp, req, db) {
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access Unauthorized")
			return
		}
		token, _ := utils.GetSessionToken(req)
		_, userName, _ := utils.GetUsernameByToken(token, db)
		err = fetchUserCreatedPosts(&posts, userName, db)

	case req.URL.Path == "/api/posts/liked":
		if !auth.SessionCheck(resp, req, db) {
			models.SendErrorResponse(resp, http.StatusUnauthorized, "Access Unauthorized")
			return
		}
		token, _ := utils.GetSessionToken(req)
		userID, _, _ := utils.GetUsernameByToken(token, db)
		err = fetchUserLikedPosts(&posts, userID, db)

	case req.URL.Path == "/api/posts" || strings.HasPrefix(req.URL.Path, "/api/posts/"):
		pathParts := strings.Split(req.URL.Path, "/")
		if len(pathParts) == 4 {
			if p, err := strconv.Atoi(pathParts[3]); err == nil && p > 0 {
				page = p
			} else {
				models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
				return
			}

		} else if len(pathParts) > 4 {
			models.SendErrorResponse(resp, http.StatusNotFound, "Page Not Found")
			return
		}
		config.Logger.Printf("Fetching posts for page %d", page)
		err = fetchAllPosts(&posts, page, db)

	default:
		models.SendErrorResponse(resp, http.StatusNotFound, "Endpoint not found")
		return
	}

	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		config.Logger.Printf("Error fetching posts: %v", err)
		return
	}

	if len(posts) == 0 {
		resp.Header().Set("Content-Type", "application/json")
		resp.Write([]byte("[]"))
		return
	}

	postJSON, err := json.Marshal(posts)
	if err != nil {
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error")
		config.Logger.Printf("Error marshaling posts: %v", err)
		return
	}

	config.Logger.Printf("Number of posts returned: %d", len(posts))
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(postJSON)
}

func fetchPostsByCategories(categories []string, posts *[]models.Post, page int, db *sql.DB) error {
	var limit = 10 * page

	var likeConditions []string
	for _, category := range categories {
		config.Logger.Printf("Category Found: %s", category)
		likeConditions = append(likeConditions, "Categories LIKE ?")
	}

	condition := strings.Join(likeConditions, " OR ")

	query := fmt.Sprintf(`
		SELECT username, id, title, content, Image_Content, Categories, created_at, likes, dislikes 
		FROM posts 
		WHERE %s
		ORDER BY created_at DESC 
		LIMIT ?`, condition)

	params := make([]any, len(categories)+2)
	for i, category := range categories {
		params[i] = "%" + category + "%"
	}
	params[len(categories)] = limit

	config.Logger.Println("Params: ", params)

	config.Logger.Printf("Executing query: %s with categories: %v, limit: %d", query, categories, limit)

	rows, err := db.Query(query, params...)
	if err != nil {
		log.Printf("Error querying posts by category: %v", err)
		return fmt.Errorf("error querying posts by category: %v", err)
	}
	defer rows.Close()

	return rowsProcess(rows, posts, db)
}

func fetchAllPosts(posts *[]models.Post, page int, db *sql.DB) error {
	var limit = 10 * page

	query := `SELECT username, id, title, content, Image_Content, Categories, created_at, likes, dislikes 
            FROM posts 
            ORDER BY created_at DESC 
            LIMIT ?`

	config.Logger.Printf("Executing query: %s with limit %d", query, limit)

	rows, err := db.Query(query, limit)
	if err != nil {
		config.Logger.Printf("Error querying all posts: %v", err)
		return fmt.Errorf("error querying all posts: %v", err)
	}
	defer rows.Close()

	return rowsProcess(rows, posts, db)
}

func rowsProcess(rows *sql.Rows, posts *[]models.Post, db *sql.DB) error {
	config.Logger.Printf("Starting to process rows...")

	for rows.Next() {
		var post models.Post

		config.Logger.Printf("Scanning row for post...")
		var creationTime time.Time
		if err := rows.Scan(&post.Username, &post.ID, &post.Title, &post.Content, &post.ImageContent, &post.Categories, &creationTime, &post.Likes, &post.Dislikes); err != nil {
			log.Printf("Error scanning post: %v", err)
			return fmt.Errorf("error scanning post: %v", err)
		}
		post.CreatedAt = utils.TimeAgo(creationTime)
		config.Logger.Printf("Successfully scanned post with ID: %d", post.ID)

		var err error
		config.Logger.Printf("Fetching comments for post ID: %d", post.ID)
		post.Comments, err = comments.GetComments(post.ID, db)

		if err != nil {
			log.Printf("Error fetching comments for post ID %d: %v", post.ID, err)
			return fmt.Errorf("error fetching comments for post ID %d: %v", post.ID, err)
		}
		config.Logger.Printf("Successfully fetched comments for post ID: %d, Comments: %v", post.ID, post.Comments)

		config.Logger.Printf("Appending post with ID: %d to posts slice", post.ID)
		*posts = append(*posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return err
	}

	config.Logger.Printf("Finished processing rows. Total posts fetched: %d", len(*posts))
	return nil
}
