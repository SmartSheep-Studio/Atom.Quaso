package models

type Subscription struct {
	Model

	AccountID  uint `json:"account_id"`
	ProviderID uint `json:"provider_id"`
}
