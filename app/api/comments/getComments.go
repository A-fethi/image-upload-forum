package comments

import (
	"database/sql"
	"forum/app/models"
)

func GetComments(postID int, db *sql.DB) ([]models.Comment, error) {
	stmt, err := db.Prepare("SELECT id, author, content, created_at, likes, dislikes FROM comments WHERE post_id = ? ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentsList []models.Comment
	for rows.Next() {
		var comment models.Comment
		if err := rows.Scan(&comment.ID, &comment.Username, &comment.Content, &comment.CreatedAt, &comment.Likes, &comment.Dislikes); err != nil {
			return nil, err
		}
		commentsList = append(commentsList, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return commentsList, nil
}

