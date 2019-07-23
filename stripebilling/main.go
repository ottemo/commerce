package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/go-ini/ini"
	"strings"

	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"

	"net/http"
	"net/smtp"
	"fmt"
	"os"
)

var (
	MailUser string
	MailPassword string
	MailServer string
	MailPort string
	MailFrom string
	MailBody string

	HttpHost string
	HttpPort string
)

func configValue(config *ini.File , key string, otherwise string) string {
	if value := strings.TrimSpace(config.Section("").Key("mail_user").String()); value != "" {
		return value
	}
	if otherwise == "{blank}" {
		panic(fmt.Sprintf("configuration value '%s' is blank ", key))
	} else {
		return otherwise
	}
}

// application start point
func main() {
	config, err := ini.Load("stripebilling.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	MailUser = configValue(config,"mail_user", "")
	MailPassword = configValue(config,"mail_password", "")
	MailServer = configValue(config,"mail_server", "127.0.0.1")
	MailPort = configValue(config,"mail_port", "25")
	MailFrom = configValue(config,"mail_from", "Stripe Billing")
	MailBody = configValue(config,"mail_body", "{blank}")

	router := httprouter.New()

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, params interface{}) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
	}

	router.GET("/", requestHandler)
	http.ListenAndServe(HttpHost, HttpPort)
}

// request handler
func requestHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

}



func SendMail(to string, subject string, body string) error {
	userName := utils.InterfaceToString(MailUser)
	password := utils.InterfaceToString(MailPassword)

	mailServer := utils.InterfaceToString(MailServer)
	mailPort := utils.InterfaceToString(MailPort)
	if mailPort != "" && mailPort != "0" {
		mailPort = ":" + mailPort
	} else {
		return nil
	}

	context := map[string]interface{}{
		"From":      MailFrom,
		"To":        to,
		"Subject":   subject,
		"Body":      body,
	}

	template := `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
Content-Type: text/html

{{.Body}}`


	body, err := utils.TextTemplate(template, context)

	var auth smtp.Auth
	if userName != "" {
		auth = smtp.PlainAuth("", userName, password, mailServer)
	}

	err = smtp.SendMail(mailServer+mailPort, auth, userName, []string{to}, []byte(body))
	return env.ErrorDispatch(err)
}
