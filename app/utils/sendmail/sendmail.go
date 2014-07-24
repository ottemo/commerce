package sendmail

import (
	"net/smtp"
	"bytes"
	"text/template"
	"strconv"
	"log"
	"fmt"
)

type EmailUser struct {
	Username    	string
	FromUserName 	string
	Password    	string
	EmailServer 	string
	Port        	int
}

type SmtpTemplateData struct {
	From    	string
	To      	string
	Subject 	string
	Body    	string
}

const emailGeneralTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}

Sincerely,

{{.From}}
`

func SendMail(to string, subject string, body string) (string, error) {
	emailUser := &EmailUser{
		Username:     "zhenyadevelop@gmail.com",
		FromUserName: "Ottemo Team",
		Password:     "111111111bb",
		EmailServer:  "smtp.gmail.com",
		Port:         587,
	}

	auth := smtp.PlainAuth("", emailUser.Username, emailUser.Password, emailUser.EmailServer)

	context := &SmtpTemplateData{
		From:    emailUser.FromUserName,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	emailTemplate := template.New("emailTemplate")
	emailTemplate, err := emailTemplate.Parse(emailGeneralTemplate)

	if err != nil {
		log.Println("error trying to parse mail template")
	}

	var doc bytes.Buffer

	err = emailTemplate.Execute(&doc, context)
	if err != nil {
		log.Println("error trying to execute mail template")
	}

	err = smtp.SendMail(
		emailUser.EmailServer + ":" + strconv.Itoa(emailUser.Port),
		auth,
		emailUser.Username,
		[]string{context.To},
		doc.Bytes() )
	if err != nil {
		log.Println("ERROR: attempting to send a mail ", err)
	}

	return fmt.Sprintf("sent status"), err;
}
