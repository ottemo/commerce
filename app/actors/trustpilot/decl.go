// Package trustpilot implements trust pilot functions
package trustpilot

import (
	"github.com/ottemo/foundation/env"
)

// Package global constants
const (
	ConstErrorModule = "trustpilot"
	ConstErrorLevel  = env.ConstErrorLevelActor

	ConstConfigPathTrustPilot               = "general.trustpilot"
	ConstConfigPathTrustPilotEnabled        = "general.trustpilot.enabled"
	ConstConfigPathTrustPilotAPIKey         = "general.trustpilot.apiKey"
	ConstConfigPathTrustPilotAPISecret      = "general.trustpilot.apiSecret"
	ConstConfigPathTrustPilotBusinessUnitID = "general.trustpilot.businessUnitID"
	ConstConfigPathTrustPilotUsername       = "general.trustpilot.username"
	ConstConfigPathTrustPilotPassword       = "general.trustpilot.password"
)

// Token holds information about last token taken and it's expiration time in Unix
type Token struct {
	Access     string
	Refresh    string
	Expiration int64
}
