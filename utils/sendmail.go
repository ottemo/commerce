package utils

import (
	"bytes"
	"net/smtp"
	"text/template"

	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
)

func SendMail(to string, subject string, body string) error {

	userName := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_USER))
	password := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_PASSWORD))

	mailServer := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_SERVER))
	mailPort := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_PORT))
	if mailPort != "" {
		mailPort = ":" + mailPort
	} else {
		mailPort = ":25"
	}

	context := map[string]interface{}{
		"From":      InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_FROM)),
		"To":        to,
		"Subject":   subject,
		"Body":      body,
		"Signature": InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_MAIL_SIGNATURE)),
	}

	emailTemplateBody := `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

{{.Signature}}`

	emailTemplate := template.New("emailTemplate")
	emailTemplate, err := emailTemplate.Parse(emailTemplateBody)
	if err != nil {
		return err
	}

	var doc bytes.Buffer
	err = emailTemplate.Execute(&doc, context)
	if err != nil {
		return err
	}

	var auth smtp.Auth = nil
	if userName != "" {
		auth = smtp.PlainAuth("", userName, password, mailServer)
	}

	err = smtp.SendMail(mailServer+mailPort, auth, userName, []string{to}, doc.Bytes())
	return err
}
