package usecase

import (
	"bytes"
	"embed"
	"html/template"
	"our-little-chatik/internal/mailer/internal"
)

//go:embed "templates"
var templateFS embed.FS

type DefaultMailUsecase struct {
	c            internal.SMTPConnector
	templateFile string
}

func NewDefaultMailUsecase(c internal.SMTPConnector,
	templateFile string) DefaultMailUsecase {
	return DefaultMailUsecase{
		c:            c,
		templateFile: templateFile,
	}
}

func (u DefaultMailUsecase) SendMailMessage(recipient string, data interface{}) error {
	// Use the ParseFS() method to parse the required template file from the embedded // file system.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+u.templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject", passing in the dynamic data and storing the // result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	// Follow the same pattern to execute the "plainBody" template and store the result // in the plainBody variable.
	plainBody := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	// And likewise with the "htmlBody" template.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	return u.c.Send(recipient, subject.String(), plainBody.String(), htmlBody.String())
}
