package models

type Notification struct {
	From    BasicUser `json:"from"`
	To      BasicUser `json:"to"`
	Message string    `json:"message"`
}
