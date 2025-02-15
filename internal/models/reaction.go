package models

type Reaction struct {
	Model
	Type      string   `json:"type"`
	UserId    string   `json:"user_id"`
	PostId    *string  `json:"post_id"`
	CommentId *string  `json:"comment_id"`
	StoryId   *string  `json:"story_id"`
	Story     *Story   `gorm:"foreignKey:StoryId" json:"story"`
	Post      *Post    `gorm:"foreignKey:PostId" json:"post"`
	Comment   *Comment `gorm:"foreignKey:CommentId" json:"comment"`
	User      User     `gorm:"foreignKey:UserId" json:"user"`
}
