package repo

import (
	"github.com/go-mail/mail/v2"
	"time"
)

// Define a Mailer struct which contains a mail.Dialer instance (used to connect to a
// SMTP server) and the sender information for your emails (the name and address you
// want the email to be from, such as "Alice Smith <alice@example.com>").
type DefaultSMTPConnector struct {
	dialer *mail.Dialer
	sender string
}

func NewDefaultSMTPConnector(host string, port int,
	username, password, sender string) DefaultSMTPConnector {
	// Initialize a new mail.Dialer instance with the given SMTP server settings. We
	// also configure this to use a 5-second timeout whenever we send an email.
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return DefaultSMTPConnector{
		dialer: dialer,
		sender: sender,
	}
}

func (c DefaultSMTPConnector) Send(recipient, subject, plainBody, htmlBody string) error {
	// Use the mail.NewMessage() function to initialize a new mail.Message instance.
	// Then we use the SetHeader() method to set the email recipient, sender and subject // headers, the SetBody() method to set the plain-text body, and the AddAlternative() // method to set the HTML body. It's important to note that AddAlternative() should // always be called *after* SetBody().
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", c.sender)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainBody)
	msg.AddAlternative("text/html", htmlBody)
	// Call the DialAndSend() method on the dialer, passing in the message to send. This // opens a connection to the SMTP server, sends the message, then closes the
	// connection. If there is a timeout, it will return a "dial tcp: i/o timeout"
	// error.
	err := c.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
