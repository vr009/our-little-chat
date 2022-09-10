package usecase

import "github.com/google/uuid"

type AuthResolver struct {
}

func NewAuthResolver() *AuthResolver {
	return &AuthResolver{}
}

func (ar *AuthResolver) ResolveToken(token string) (uuid.UUID, error) {
	id, err := uuid.Parse(token)
	return id, err
}
