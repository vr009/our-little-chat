package models

type NotificationType string

const (
	InfoMessage NotificationType = "info"
	ChatMessage NotificationType = "chat"
)

// Notification is a type that gets encoded into a json document when communicating
// with peer service. Type field describes whether the notification has service info
// purpose or it is a message. Body may contain models.Message or other service structures.
type Notification struct {
	Type NotificationType `json:"type,omitempty"`
	Body interface{}      `json:"body,omitempty"`
}

type ConnectionStatusType string

const (
	Established ConnectionStatusType = "established"
	Failed      ConnectionStatusType = "failed"
	Conflict    ConnectionStatusType = "conflict"
)

// PeerConnectionStatus is a type for notifying peer about a connection status.
type PeerConnectionStatus struct {
	Status     ConnectionStatusType
	Properties map[string]any
}
