package session

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/url"
	"sync"
	"time"

	"errors"
)

const (
	ALPHANUMERIC = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

	SESSION_COOKIE_NAME = "OTTEMOSESSION"
)

var (
	Sessions      = make(map[string]*Session)
	sessionsMutex sync.RWMutex

	gcRate int64 = 10
)

// returns session object for given id or nil
func GetSessionById(sessionId string) (*Session, error) {
	if session, ok := Sessions[sessionId]; ok == true {
		return session, nil
	} else {
		return nil, errors.New("session not found")
	}
}

// returns session object for request or creates new one
func StartSession(request *http.Request, responseWriter http.ResponseWriter) (*Session, error) {

	// check session-cookie
	cookie, err := request.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			return newSession(responseWriter)
		}
		return nil, err
	}

	// session-cookie found
	sessionId := cookie.Value
	if session, ok := Sessions[sessionId]; ok == true {
		return session, nil
	}

	return newSession(responseWriter)
}

// initializes new session
func newSession(responseWriter http.ResponseWriter) (*Session, error) {

	// receiving new session id
	sessionId, err := newSessionId()
	if err != nil {
		return nil, err
	}

	// initializing session structure
	sessionId = url.QueryEscape(sessionId)
	Sessions[sessionId] = &Session{
		id:     sessionId,
		values: make(map[string]interface{}),
		time:   time.Now()}

	// updating cookies
	cookie := &http.Cookie{Name: SESSION_COOKIE_NAME, Value: sessionId, Path: "/"}
	http.SetCookie(responseWriter, cookie)

	// garbage collecting
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(gcRate))
	if err == nil && randomNumber.Cmp(big.NewInt(1)) == 0 {
		Gc()
	}

	return Sessions[sessionId], nil
}

// returns new session number
func newSessionId() (string, error) {
	sessionId := make([]byte, 32)
	if _, err := rand.Read(sessionId); err != nil {
		return "", errors.New("can't generate sessionId")
	}

	for i := 0; i < 32; i++ {
		sessionId[i] = ALPHANUMERIC[sessionId[i]%62]
	}

	return string(sessionId), nil
}

// removes expired sessions
func Gc() {
	for id, session := range Sessions {
		if time.Now().Sub(session.time).Seconds() > 3600 {
			sessionsMutex.Lock()

			delete(Sessions, id)

			sessionsMutex.Unlock()
		}
	}
}