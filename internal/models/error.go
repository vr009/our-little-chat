package models

type StatusCode int

const (
	OK StatusCode = iota
	BadRequest
	NotFound
	Forbidden
	InternalError
	Deleted
	Conflict
)

type Error struct {
	Msg string `json:"message"`
}
