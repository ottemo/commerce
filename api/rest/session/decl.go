package session

import (
	"net/http"
	"crypto/rand"
	"time"
	"sync"
	"net/url"
	"math/big"

	"errors"
)

var (
	CookieName = "OTTEMOSESSION"
	Sessions map[string]*Session = make(map[string]*Session)

	gcRate int64 = 10
	sessionsMutex sync.RWMutex
)

// returns session object for request or creates new one
func StartSession(request *http.Request, responseWriter http.ResponseWriter) (*Session, error){

	// check session-cookie
	cookie, err := request.Cookie(CookieName)
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
		time:   time.Now() }


	// updating cookies
	cookie := &http.Cookie{Name: CookieName, Value: sessionId, Path: "/"}
http.SetCookie(responseWriter, cookie)


// garbage collecting
randomNumber, err := rand.Int(rand.Reader, big.NewInt(gcRate))
if err == nil && randomNumber.Cmp(big.NewInt(1)) == 0 {
Gc()
}

return Sessions[sessionId], nil
}


// returns new session number
func newSessionId() (string, error){
	sessionId := make([]byte, 32)
	if _, err := rand.Read(sessionId); err != nil {
		return "", errors.New("can't generate sessionId")
	}
	return string(sessionId), nil
}


// removes expired sessions
func Gc() {
	for id, session := range Sessions {
		if time.Now().Sub(session.time).Seconds() > 3600 {
			sessionsMutex.Lock()

			Sessions[id] = nil

			sessionsMutex.Unlock()
		}
	}
}
