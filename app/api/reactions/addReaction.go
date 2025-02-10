package reactions

import (
	"database/sql"
	"encoding/json"
	"forum/app/api/auth"
	"forum/app/config"
	"forum/app/models"
	"forum/app/utils"
	"net/http"
	"strconv"
)

func AddReaction(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	config.Logger.Println("Received request to add reaction")

	if req.Method != http.MethodPost {
		config.Logger.Println("Invalid HTTP method:", req.Method)
		models.SendErrorResponse(resp, http.StatusMethodNotAllowed, "Method Not Alowed")
		return
	}

	if !auth.SessionCheck(resp, req, db) {
		config.Logger.Println("Unauthorized access attempt")
		models.SendErrorResponse(resp, http.StatusUnauthorized, "Access: User not authenticated, Login and try!!")
		return
	}

	token, _ := utils.GetSessionToken(req)
	config.Logger.Println("Retrieved session token:", token)

	userID, _, err := utils.GetUsernameByToken(token, db)
	if err != nil {
		config.Logger.Printf("Error getting username by token: %v", err)
		models.SendErrorResponse(resp, 500, "Error: Internal Server Error, try again later!")
		return
	}
	config.Logger.Println("Authenticated user ID:", userID)

	var reactRequest struct {
		ItemID   string `json:"item_id"`
		ItemType string `json:"item_type"`
		Action   string `json:"action"`
	}

	if err := json.NewDecoder(req.Body).Decode(&reactRequest); err != nil {
		config.Logger.Printf("Error decoding request body: %v", err)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid Request Body")
		return
	}

	config.Logger.Printf("Processing reaction request: %+v", reactRequest)

	itemID, err := strconv.Atoi(reactRequest.ItemID)
	if err != nil || (reactRequest.Action != "like" && reactRequest.Action != "dislike") {
		config.Logger.Println("Invalid item ID or action")
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: Invalid reaction")
		return
	}

	tableName := ""
	switch reactRequest.ItemType {
	case "post":
		tableName = "posts"
	case "comment":
		tableName = "comments"
	default:
		config.Logger.Println("Invalid item type:", reactRequest.ItemType)
		models.SendErrorResponse(resp, http.StatusBadRequest, "Error: item Not a Post/Comment")
		return
	}

	config.Logger.Printf("Target table: %s, Item ID: %d, Action: %s", tableName, itemID, reactRequest.Action)

	var existingID int
	var existingAction string
	query := `SELECT id, action FROM user_interactions WHERE user_id = ? AND item_id = ? AND item_type = ?`
	err = db.QueryRow(query, userID, itemID, reactRequest.ItemType).Scan(&existingID, &existingAction)

	if err == sql.ErrNoRows {
		config.Logger.Println("No existing interaction, adding new entry")

		_, err = db.Exec(
			`INSERT INTO user_interactions (user_id, item_id, item_type, action) VALUES (?, ?, ?, ?)`,
			userID, itemID, reactRequest.ItemType, reactRequest.Action,
		)
		if err != nil {
			config.Logger.Printf("Error adding interaction: %v", err)
			models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error.")
			return
		}

		updateColumn := "likes"
		if reactRequest.Action == "dislike" {
			updateColumn = "dislikes"
		}
		config.Logger.Printf("Incrementing %s count for item ID: %d", updateColumn, itemID)
		db.Exec(`UPDATE `+tableName+` SET `+updateColumn+` = `+updateColumn+` + 1 WHERE id = ?`, itemID)

	} else if err == nil {
		config.Logger.Printf("Existing interaction found (ID: %d), current action: %s", existingID, existingAction)

		if existingAction == reactRequest.Action {
			config.Logger.Println("User is removing their existing reaction")
			updateColumn := "likes"
			if reactRequest.Action == "dislike" {
				updateColumn = "dislikes"
			}
			db.Exec(`UPDATE `+tableName+` SET `+updateColumn+` = CASE WHEN `+updateColumn+` > 0 THEN `+updateColumn+` - 1 ELSE 0 END WHERE id = ?`, itemID)
			db.Exec(`DELETE FROM user_interactions WHERE id = ?`, existingID)

		} else {
			config.Logger.Println("User is changing their reaction")
			db.Exec(`DELETE FROM user_interactions WHERE id = ?`, existingID)
			_, err = db.Exec(
				`INSERT INTO user_interactions (user_id, item_id, item_type, action) VALUES (?, ?, ?, ?)`,
				userID, itemID, reactRequest.ItemType, reactRequest.Action,
			)
			if reactRequest.Action == "like" {
				db.Exec(`UPDATE `+tableName+` SET likes = likes + 1, dislikes = CASE WHEN dislikes > 0 THEN dislikes - 1 ELSE 0 END WHERE id = ?`, itemID)
			} else {
				db.Exec(`UPDATE `+tableName+` SET dislikes = dislikes + 1, likes = CASE WHEN likes > 0 THEN likes - 1 ELSE 0 END WHERE id = ?`, itemID)
			}
		}
	} else {
		config.Logger.Printf("Error querying interactions: %v", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error.")
		return
	}

	var updatedLikes, updatedDislikes int
	err = db.QueryRow(`SELECT likes, dislikes FROM `+tableName+` WHERE id = ?`, itemID).Scan(&updatedLikes, &updatedDislikes)
	if err != nil {
		config.Logger.Printf("Error retrieving updated likes/dislikes: %v", err)
		models.SendErrorResponse(resp, http.StatusInternalServerError, "Error: Internal Server Error.")
		return
	}

	config.Logger.Printf("Updated counts - Likes: %d, Dislikes: %d", updatedLikes, updatedDislikes)

	response := map[string]interface{}{
		"item_id":  itemID,
		"likes":    updatedLikes,
		"dislikes": updatedDislikes,
	}
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(response)
	config.Logger.Println("Response sent successfully")
}
