package models

import (
	"gorm.io/datatypes"
	"time"
)

type Post struct {
	Model

	Content     string                      `json:"content"`
	Tags        datatypes.JSONSlice[string] `json:"tags"`
	Attachments datatypes.JSONSlice[string] `json:"attachments"`
	Type        string                      `json:"type"`
	Scope       string                      `json:"scope"`
	IpAddress   string                      `json:"ip_address"`
	IsEdited    bool                        `json:"is_edited"`
	IsHidden    bool                        `json:"is_hidden"`
	Comments    []Post                      `json:"comments" gorm:"foreignKey:BelongID"`
	Likes       []Like                      `json:"likes" gorm:"foreignKey:PostID"`
	Dislikes    []Dislike                   `json:"dislikes" gorm:"foreignKey:PostID"`
	PublishedAt time.Time                   `json:"published_at"`
	BelongID    *uint                       `json:"belong_id"`
	AccountID   uint                        `json:"account_id"`
}
