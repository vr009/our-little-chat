package internal

type SMTPConnector interface {
	Send(recipient, subject, plainBody, htmlBody string) error
}

type MailerUsecase interface {
	SendMailMessage(recipient string, data interface{}) error
}
