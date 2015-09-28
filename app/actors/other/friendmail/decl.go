// Package friendmail is an extension which provides ability to send email to friend
package friendmail

import (
	"github.com/ottemo/foundation/env"
	"time"
	"sync"
	"github.com/dchest/captcha"
)

// Package global constants
const (
	ConstCollectionNameFriendMail = "friend_mail"

	ConstConfigPathGroup   = "general.friendmail"
	ConstConfigPathEmailTemplate = "general.friendmail.template"
	ConstConfigPathEmailSubject = "general.friendmail.subject"

	ConstErrorModule = "friendmail"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstMaxCaptchaItems = 100000 // captcha maximum amount (to prevent memory leaks)
	ConstCaptchaLifeTime = 300    // seconds generated captcha works (5 min)
)


var (
	captchaValuesMutex sync.RWMutex      // synchronization on captchaValues variable
	captchaValues map[string]time.Time   // global variable to track generated captcha codes
	captchaStore captcha.Store
)