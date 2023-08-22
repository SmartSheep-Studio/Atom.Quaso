package models

type Account struct {
	Model

	Posts         []Post         `json:"posts"`
	LikedPosts    []Like         `json:"liked_posts" gorm:"foreignKey:AccountID"`
	DislikedPosts []Dislike      `json:"disliked_posts" gorm:"foreignKey:AccountID"`
	Subscriptions []Subscription `json:"subscriptions" gorm:"foreignKey:AccountID"`
	Subscribers   []Subscription `json:"subscribers" gorm:"foreignKey:ProviderID"`
	UserID        uint           `json:"user_id"`
}
