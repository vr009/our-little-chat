package models

type UpdateOptions struct {
	Action string
}

const (
	UpdatePhotoURL              = "update_photo_url"
	AddUsersToParticipants      = "add_users_to_participants"
	RemoveUsersFromParticipants = "remove_user_from_participants"
)
