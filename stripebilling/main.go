package main

/*
Generate private key (.key)

Key considerations for algorithm "RSA" ≥ 2048-bit
> openssl genrsa -out server.key 2048

Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
> openssl ecparam -genkey -name secp384r1 -out server.key

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
> openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

*/

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/go-ini/ini"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ottemo/commerce/env"
	"github.com/ottemo/commerce/utils"

	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

const (
	alphanumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	tokenLifetime = time.Minute * 30;
)

var (
	MailUser     string
	MailPassword string
	MailHost     string
	MailPort     string
	MailFrom     string
	MailSubject  string
	MailFile     string
	MailTemplate string

	JenkinsUrl string
	JenkinsStatusUrl string
	JenkinsUser string

	HttpHost string
	HttpPort string
	HttpCertFile string
	HttpKeyFile string

	tokens map[string]map[string]interface{}
	tokens_mutex sync.Mutex
)

// application start point
func main() {
	tokens = make(map[string]map[string]interface{})

	// configuration reading
	config, err := ini.Load("stripebilling.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	HttpCertFile = configValue(config,"http_cert_file", "")
	HttpKeyFile = configValue(config,"http_key_file", "")

	defaultHttpPort := "80"
	if HttpCertFile != "" {
		defaultHttpPort = "433"
	}

	HttpHost = configValue(config,"http_address", "127.0.0.1")
	HttpPort = configValue(config,"http_port", defaultHttpPort)

	JenkinsUrl = configValue(config,"jenkins_url", "{error}")
	JenkinsUser = configValue(config,"jenkins_user", "{error}")
	JenkinsStatusUrl = configValue(config,"jenkins_status_url", "{error}")

	MailUser = configValue(config,"mail_user", "")
	MailPassword = configValue(config,"mail_password", "")
	MailHost = configValue(config,"mail_host", "127.0.0.1")
	MailPort = configValue(config,"mail_port", "")
	MailFrom = configValue(config,"mail_from", "Stripe Billing")
	MailFrom = configValue(config,"mail_subject", "Stripe Billing")
	MailFile = configValue(config,"mail_file", "stripebilling.tpl")

	if MailPort != "" && MailPort != "0" {
		MailPort = ":" + MailPort
	}

	// mail template file reading
	file, err := os.Open(MailFile)
	if err != nil {
		panic(err)
	}
	template, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	MailTemplate = string(template)

	// tokens cleaner
	ticker := time.NewTicker(tokenLifetime + 10)
	go func() {
		for _ = range ticker.C {
			now := time.Now()
			log("tokens cleanup started at %v", now)
			tokens_mutex.Lock()
			for token, context := range tokens {
				if value, present := context["Expire"]; present {
					if expire, ok := value.(time.Time); !ok || expire.Sub(now) < 0 {
						delete(tokens, token)
						fmt.Printf("token %s expired\n", token)
					}
				} else {
					delete(tokens, token)
				}
			}
			tokens_mutex.Unlock()
			log("tokens cleanup finished at %v", now)
		}
	}()


	// http server setup
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, params interface{}) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("page not found"))
	}
	router.GET("/:token", requestHandler)
	router.POST("/", requestHandler)

	if HttpCertFile != "" && HttpKeyFile != "" {
		log("Starting HTTP/TLS server on %s:%s\n", HttpHost, HttpPort)
		err := http.ListenAndServeTLS(fmt.Sprintf("%s:%s", HttpHost, HttpPort), HttpCertFile, HttpKeyFile, router)
		if err != nil {
			panic(err)
		}
	} else {
		log("Starting HTTP server on %s:%s\n", HttpHost, HttpPort)
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", HttpHost, HttpPort), router)
		if err != nil {
			panic(err)
		}
	}
}

// performs log operation
func log(format string, a ...interface{}) {
	if format[len(format)-1:] != "\n" {
		format = format + "\n"
	}
	format = time.Now().Format(time.RFC3339) + ":" + format
	fmt.Printf(format, a...)
}

// GenerateSessionID returns new session id number
func getToken() string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}

	for i := 0; i < 32; i++ {
		token[i] = alphanumeric[token[i]%62]
	}

	return string(token)
}

// takes the config value from ini file
func configValue(config *ini.File , key string, otherwise string) string {
	if value := strings.TrimSpace(config.Section("").Key(key).String()); value != "" {
		return value
	}
	if otherwise == "{error}" {
		panic(fmt.Sprintf("configuration value '%s' is blank ", key))
	} else {
		return otherwise
	}
}

// writes error to http writer
func writeError(writer http.ResponseWriter, err interface{}) {
	message := fmt.Sprintf("Error: %s", err)
	log("Error: %s", message)
	writer.WriteHeader(400)
	writer.Write([]byte(message))
}

// request handler
func requestHandler(response http.ResponseWriter, request *http.Request, params httprouter.Params) {
	log("Handling %s request", request.Method)

	contentType := request.Header.Get("Content-Type")
	context := map[string]interface{}{
		"From":    MailFrom,
		"To":      "",
		"Subject": MailSubject,
		"Body":    MailFile,
	}

	if request.Method == "POST" {
		if strings.Contains(contentType, "form-data") {
			if err := request.ParseForm(); err != nil {
				writeError(response, err);
			}
			for attribute, value := range request.PostForm {
				context[attribute], _ = url.QueryUnescape(value[0])
			}

			if err := request.ParseMultipartForm(32 << 20); err != nil { // 32 MB
				writeError(response, err);
			}
			if request.MultipartForm != nil {
				for attribute, value := range request.MultipartForm.Value {
					context[attribute], _ = url.QueryUnescape(value[0])
				}
			}
		} else if strings.Contains(contentType, "urlencode") {
			if err := request.ParseForm(); err != nil {
				writeError(response, err);
			}
			for attribute, value := range request.PostForm {
				context[attribute], _ = url.QueryUnescape(value[0])
			}
		} else {
			data := make(map[string]interface{})
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				writeError(response, err);
			}
			if err := json.Unmarshal(body, &data); err != nil {
				writeError(response, err);
			}
			for key, value := range data {
				context[key] = value
			}
		}
	}

	if token:= params.ByName("token"); token != "" {
		tokens_mutex.Lock()
		data, present := tokens[token]
		tokens_mutex.Unlock()

		if present {
			log("processing token %s", tokens)

			for key, value := range data {
				if _, present := context[key]; !present {
					context[key] = value
				}
			}
			jenkinsCall(response, context)
		} else {
			writeError(response, "Invalid token: " + token)
		}
	} else {
		if to := utils.GetFirstMapValue(context, "email", "e-mail", "to"); to != nil {
			token := getToken()
			context["Token"] = token
			context["To"] = to
			context["Expire"] = tokenLifetime

			log("token %s created", tokens)

			if err := sendMail(utils.InterfaceToString(to), context); err != nil {
				writeError(response, err)
			}

			tokens_mutex.Lock()
			tokens[token] = context
			tokens_mutex.Unlock()
		} else {
			writeError(response, "Recipient was not found, expecting 'email' or 'to' field: " + token)
		}
	}
}

func jenkinsCall(response http.ResponseWriter, context map[string]interface{}) {
	var requestUrl string

	if _, present := context["Started"]; present {
		log("Requesting status for a job started at %v, token %s", context["Started"], context["Token"])
		requestUrl = JenkinsStatusUrl
	} else {
		log("Performing Jenkins call for token %s", context["Token"])
		requestUrl = JenkinsUrl
		context["Started"] = time.Now()
	}

	data := make(url.Values)
	for key, value := range context {
		data[key] = []string{utils.InterfaceToString(value)}
	}

	client := new(http.Client)

	request, err := http.NewRequest("POST", requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		writeError(response, err)
	}
	request.Header.Set("Authorization", "Basic "+ base64.StdEncoding.EncodeToString([]byte(JenkinsUser)))

	answer, err := client.Do(request)
	if err != nil{
		writeError(response, err)
	}

	body, err := ioutil.ReadAll(answer.Body)
	if err != nil {
		writeError(response,err)
	}

	response.Header().Set("Content-Type", answer.Header.Get("Content-Type"))
	response.WriteHeader(answer.StatusCode)
	response.Write(body)
}

// sends email
func sendMail(to string, context map[string]interface{}) error {
	log("Sending e-mail to '%s', token '%s'", to, context["Token"])

	body, err := utils.TextTemplate(MailTemplate, context)
	if err != nil {
		log("SendMail error: %s", err.Error())
	}

	var auth smtp.Auth
	if MailUser != "" {
		auth = smtp.PlainAuth("", MailUser, MailPassword, MailHost)
	}
	err = smtp.SendMail(MailHost+MailPort, auth, MailFrom, []string{to}, []byte(body))
	return env.ErrorDispatch(err)
}
