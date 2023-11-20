package models

type ActivationTask struct {
	ActivationCode string `json:"activation_code"`
	Receiver       string `json:"receiver"`
}
