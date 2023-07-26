package models

type Metadata struct {
	ID          string `bson:"id" json:"id" `
	Image       string `json:"image"`
	Description string `json:"description"`
	Name        string `json:"name"`
}
