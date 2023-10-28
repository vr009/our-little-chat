package models

import (
	"github.com/google/uuid"
	"our-little-chatik/internal/pkg/validator"
)

type CreateChatRequest struct {
	Participants []uuid.UUID `json:"participants,omitempty"`
	Name         *string     `json:"name,omitempty"`
	PhotoURL     *string     `json:"photo_url,omitempty"`
	IssuerID     uuid.UUID   `json:"-"`
}

func ValidateCreateChatRequest(v *validator.Validator, request CreateChatRequest) {
	if request.Name != nil {
		v.Check(len(*request.Name) < 50, "name", "must be less than 50 bytes")
	}
	if request.PhotoURL != nil {
		//TODO add check for extensions??
	}
	v.Check(request.Participants != nil, "participants", "must be provided")
	if request.Participants != nil {
		v.Check(len(request.Participants) < 100, "participants", "must be less than 100 members")
	}
}

type AddUsersToChatRequest struct {
	ChatID       *uuid.UUID  `json:"chat_id"`
	Participants []uuid.UUID `json:"participants,omitempty"`
}

func ValidateAddUsersToChatRequest(v *validator.Validator, request AddUsersToChatRequest) {
	v.Check(request.ChatID != nil, "chat_id", "must be provided")
	if request.ChatID != nil {
		v.Check(*request.ChatID != uuid.Nil, "chat_id", "must be a correct uuid value")
	}
	v.Check(request.Participants != nil, "participants", "must be provided")
	v.Check(len(request.Participants) < 10, "participants", "can't add more than 10 users")
}

type RemoveUsersFromChatRequest struct {
	ChatID       *uuid.UUID  `json:"chat_id"`
	Participants []uuid.UUID `json:"participants,omitempty"`
}

func ValidateRemoveUsersFromChatRequest(v *validator.Validator, request RemoveUsersFromChatRequest) {
	v.Check(request.ChatID != nil, "chat_id", "must be provided")
	if request.ChatID != nil {
		v.Check(*request.ChatID != uuid.Nil, "chat_id", "must be a correct uuid value")
	}
	v.Check(request.Participants != nil, "participants", "must be provided")
	v.Check(len(request.Participants) < 10, "participants", "can't remove more than 10 users")
}

type UpdateChatPhotoURLRequest struct {
	ChatID   *uuid.UUID `json:"chat_id"`
	PhotoURL *string    `json:"photo_url,omitempty"`
}

func ValidateUpdateChatPhotoURLRequest(v *validator.Validator, request UpdateChatPhotoURLRequest) {
	v.Check(request.ChatID != nil, "chat_id", "must be provided")
	if request.ChatID != nil {
		v.Check(*request.ChatID != uuid.Nil, "chat_id", "must be a correct uuid value")
	}
	v.Check(request.PhotoURL != nil, "photo_url", "must be provided")
	if request.PhotoURL != nil {
		//TODO add check for extensions??
	}
}
