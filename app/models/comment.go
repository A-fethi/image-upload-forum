package models

type Comment struct {
	ID        int
	Content   string
	Username  string
	PostID    int
	CreatedAt string
	Likes     int 
	Dislikes  int 
}
