package models

import "our-little-chatik/internal/models"

type NotificationType int8

const (
	InfoMessage NotificationType = iota
	ChatMessage
)

type Notification struct {
	Type        NotificationType `json:"type,omitempty"`
	Message     *models.Message  `json:"message,omitempty"`
	Description string           `json:"description,omitempty"`
}
