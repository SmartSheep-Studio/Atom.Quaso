package models

type Account struct {
	Model

	Nickname    string     `json:"nickname"`
	Posts       []Post     `json:"posts"`
	Subscribers []*Account `gorm:"many2many:account_followers"`
	Followers   []*Account `gorm:"many2many:account_followers"`
	UserID      uint       `json:"user_id"`
}
