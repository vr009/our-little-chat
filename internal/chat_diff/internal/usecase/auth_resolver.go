package usecase

import "github.com/google/uuid"

type AuthResolver struct {
}

func NewAuthResolver() *AuthResolver {
	return &AuthResolver{}
}

func (ar *AuthResolver) ResolveToken(token string) (uuid.UUID, error) {
	id, err := uuid.Parse("62391bd9-157c-4513-8e7c-c082e00d2b7e")
	return id, err
}
