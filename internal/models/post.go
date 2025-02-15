package models

type Post struct {
	Model
	Content   string     `json:"content"`
	AuthorId  uint       `json:"author_id"`
	Author    User       `gorm:"foreignKey:AuthorId" json:"author"`
	Comments  []Comment  `gorm:"foreignKey:PostId" json:"comments"`
	Reactions []Reaction `gorm:"foreignKey:PostId" json:"reactions"`
	Media     []Media    `gorm:"foreignKey:PostId" json:"media"`
}
