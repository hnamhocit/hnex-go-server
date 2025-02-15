package models

type Story struct {
	Model
	Text            *string    `json:"text"`
	TextCoordinates string     `json:"text_coordinates" gorm:"type:JSON;default:'{\"x\":0,\"y\":0}'"`
	MediaId         uint       `json:"media_id" gorm:"unique"`
	Media           Media      `json:"media" gorm:"foreignKey:MediaId"`
	AuthorId        uint       `json:"author_id"`
	Author          User       `json:"author" gorm:"foreignKey:AuthorId"`
	Viewers         []User     `json:"viewers" gorm:"many2many:story_viewers;"`
	Reactions       []Reaction `json:"reactions" gorm:"foreignKey:StoryId"`
}
