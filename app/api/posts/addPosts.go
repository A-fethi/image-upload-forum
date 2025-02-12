package posts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
)

func AddPost(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var postID int

	if !auth.SessionCheck(resp, req, db) {
		http.Error(resp, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err := req.ParseMultipartForm(20 << 20)
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(resp, "Error parsing form data", http.StatusBadRequest)
		return
	}

	title := req.FormValue("title")
	postContent := req.FormValue("content")
	categories := req.Form["categories"]

	var imageFilename string
	file, header, err := req.FormFile("image")

	if err == nil {
		defer file.Close()

		if err := utils.ValidateImage(header.Filename); err != nil {
			models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid File name extention")
			return
		}
		imageFilename = fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		var imagePath string
		path := "static/uploads"
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			os.Mkdir("static/uploads", os.ModePerm)
		}
		imagePath = filepath.Join("static/uploads", imageFilename)

		outFile, err := os.Create(imagePath)
		if err != nil {
			log.Println("Error creating file:", err)
			http.Error(resp, "Error saving image", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		_, err = io.Copy(outFile, file)
		if err != nil {
			log.Println("Error saving file:", err)
			http.Error(resp, "Error saving image", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("No image uploaded, continuing without image.")
	}

	if err := utils.ValidatePost(title, postContent); err != nil {
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Title/Post")
		return
	}

	sessionToken, err := utils.GetSessionToken(req)
	if err != nil || sessionToken == "" || !auth.SessionCheck(resp, req, db) {
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: Unauthorized")
		return
	}

	catCHECK := utils.CategoriesCheck(categories)
	if !catCHECK {
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Categories")
		return
	}

	_, username, err := utils.GetUsernameByToken(sessionToken, db)
	if err != nil {
		config.Logger.Println("Failed to get username:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		return
	}

	err = db.QueryRow(`
		INSERT INTO posts (username, title, content, image_Content, categories, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING id`,
		username, title, postContent, imageFilename, strings.Join(categories, " "), time.Now()).Scan(&postID)
	if err != nil {
		config.Logger.Println("Failed to insert post:", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error. Try later")
		return
	}

	post := models.Post{
		Username:     username,
		ID:           postID,
		Title:        title,
		Content:      postContent,
		ImageContent: imageFilename,
		Categories:   strings.Join(categories, " "),
		CreatedAt:    utils.TimeAgo(time.Now()),
		Likes:        0,
		Dislikes:     0,
	}

	config.Logger.Printf("Post created successfully, postID: %d", postID)
	resp.WriteHeader(http.StatusCreated)
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(post)
}
