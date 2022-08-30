package models

type StorageConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	Database    string `json:"database"`
	MaxAttempts int    `json:"max_attempts"`
}
