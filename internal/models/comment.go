package models

type Comment struct {
	Model
	Media           []Media    `json:"media" gorm:"foreignKey:CommentId"`
	Reactions       []Reaction `json:"reactions" gorm:"foreignKey:CommentId"`
	UserId          uint       `json:"user_id"`
	User            User       `json:"user" gorm:"foreignKey:UserId"`
	PostId          uint       `json:"post_id"`
	Post            Post       `json:"post" gorm:"foreignKey:PostId"`
	Content         string     `json:"content"`
	ParentCommentId *string    `json:"parent_comment_id"`
	ParentComment   *Comment   `json:"parent_comment" gorm:"foreignKey:ParentCommentId"`
}
