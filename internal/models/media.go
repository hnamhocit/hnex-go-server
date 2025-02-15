package models

type Media struct {
	Model
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Size        int64    `json:"size"`
	ContentType string   `json:"content_type"`
	UserID      uint     `json:"user_id"`
	User        User     `gorm:"foreignKey:UserID"`
	PostId      *uint    `json:"post_id"`
	Post        *Post    `gorm:"foreignKey:PostId"`
	CommentId   *uint    `json:"comment_id"`
	Comment     *Comment `gorm:"foreignKey:CommentId"`
}
