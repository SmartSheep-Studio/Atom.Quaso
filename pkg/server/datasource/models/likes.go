package models

type Like struct {
	Model

	AccountID uint `json:"account_id"`
	PostID    uint `json:"post_id"`
}

type Dislike struct {
	Model

	AccountID uint `json:"account_id"`
	PostID    uint `json:"post_id"`
}
