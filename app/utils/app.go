package utils

import (
	"github.com/ottemo/foundation/app"
	"github.com/ottemo/foundation/env"
	"strings"
)

// returns url related to dashboard server
func GetDashboardUrl(path string) string {
	baseUrl := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_DASHBOARD_URL))
	return strings.TrimRight(baseUrl, "/") + "/" + path
}

// returns url related to storefront server
func GetStorefrontUrl(path string) string {
	baseUrl := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_STOREFRONT_URL))
	return strings.TrimRight(baseUrl, "/") + "/" + path
}

// returns url related to foundation server
func GetFoundationUrl(path string) string {
	baseUrl := InterfaceToString(env.ConfigGetValue(app.CONFIG_PATH_FOUNDATION_URL))
	return strings.TrimRight(baseUrl, "/") + "/" + path
}
