package models

import "our-little-chatik/internal/models"

type NotificationType string

const (
	InfoMessage NotificationType = "info"
	ChatMessage NotificationType = "chat"
)

type Notification struct {
	Type        NotificationType `json:"type,omitempty"`
	Message     *models.Message  `json:"message,omitempty"`
	Description string           `json:"description,omitempty"`
}
