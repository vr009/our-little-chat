package models

import (
	"github.com/google/uuid"
	"our-little-chatik/internal/pkg/validator"
)

type UpdateUserRequest struct {
	Nickname    *string `json:"nickname,omitempty"`
	Name        *string `json:"name,omitempty"`
	Surname     *string `json:"surname,omitempty"`
	OldPassword *string `json:"old_password,omitempty"`
	NewPassword *string `json:"new_password,omitempty"`
	Avatar      *string `json:"avatar,omitempty"`
}

func ValidateUpdateUserRequest(v *validator.Validator, request UpdateUserRequest) {
	if request.Name != nil {
		v.Check(len(*request.Name) < 50, "name", "must be less than 50 symbols")
	}
	if request.Nickname != nil {
		v.Check(len(*request.Nickname) < 50, "nickname", "must be less than 50 symbols")
	}
	if request.Surname != nil {
		v.Check(len(*request.Surname) < 50, "surname", "must be less than 50 symbols")
	}
	if request.NewPassword != nil {
		v.Check(request.OldPassword != nil, "old password", "must be provided and correct for update")
		v.Check(request.NewPassword != nil, "old password", "must be provided and correct")
		if request.OldPassword != nil {
			ValidatePasswordPlaintext(v, *request.OldPassword)
		}
		ValidatePasswordPlaintext(v, *request.NewPassword)
	}
	if request.OldPassword != nil && request.NewPassword == nil {
		v.AddError("new password", "must be provided for update")
	}
}

type SignUpPersonRequest struct {
	Nickname *string `json:"nickname,omitempty"`
	Name     *string `json:"name,omitempty"`
	Surname  *string `json:"surname,omitempty"`
	Password *string `json:"password,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

func ValidateSignUpRequest(v *validator.Validator, request SignUpPersonRequest) {
	v.Check(request.Name != nil, "name", "must be provided")
	v.Check(request.Surname != nil, "surname", "must be provided")
	v.Check(request.Nickname != nil, "nickname", "must be provided")
	v.Check(request.Password != nil, "password", "must be provided")
	if request.Name != nil {
		v.Check(len(*request.Name) < 500, "name", "must not be more than 500 bytes long")
	}
	if request.Nickname != nil {
		v.Check(len(*request.Nickname) < 500, "nickname", "must not be more than 500 bytes long")
	}
	if request.Surname != nil {
		v.Check(len(*request.Surname) < 500, "surname", "must not be more than 500 bytes long")
	}
	if request.Password != nil {
		ValidatePasswordPlaintext(v, *request.Password)
	}
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

type LoginRequest struct {
	Nickname *string `json:"nickname,omitempty"`
	Password *string `json:"password,omitempty"`
}

func ValidateLoginRequest(v *validator.Validator, request LoginRequest) {
	v.Check(request.Nickname != nil, "nickname", "must be provided")
	v.Check(request.Password != nil, "password", "must be provided")
	if request.Nickname != nil {
		v.Check(len(*request.Nickname) < 500, "nickname", "must not be more than 500 bytes long")
	}
	if request.Password != nil {
		ValidatePasswordPlaintext(v, *request.Password)
	}
}

type GetUserRequest struct {
	UserID uuid.UUID
}

func ValidateGetUserRequest(v *validator.Validator, request GetUserRequest) {
	v.Check(request.UserID != uuid.Nil, "UserID", "must be a correct value")
}
