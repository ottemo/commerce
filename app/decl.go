// Package app represents application object.
//
// That package contains routines to register callbacks on application start/end,
// API functions for administrator login, system configuration values, etc.
//
// In general this package contains the stuff addressed to system application.
package app

// Package global constants
const (
	CONFIG_PATH_GROUP       = "general"
	CONFIG_PATH_APP_GROUP   = "general.app"
	CONFIG_PATH_STORE_GROUP = "general.store"
	CONFIG_PATH_MAIL_GROUP  = "general.mail"

	CONFIG_PATH_STOREFRONT_URL = "general.app.storefront_url"
	CONFIG_PATH_DASHBOARD_URL  = "general.app.dashboard_url"
	CONFIG_PATH_FOUNDATION_URL = "general.app.foundation_url"

	CONFIG_PATH_STORE_NAME  = "general.store.name"
	CONFIG_PATH_STORE_EMAIL = "general.store.email"

	CONFIG_PATH_STORE_ROOT_LOGIN    = "general.store.root_login"
	CONFIG_PATH_STORE_ROOT_PASSWORD = "general.store.root_password"

	CONFIG_PATH_STORE_COUNTRY      = "general.store.country"
	CONFIG_PATH_STORE_STATE        = "general.store.state"
	CONFIG_PATH_STORE_CITY         = "general.store.city"
	CONFIG_PATH_STORE_ADDRESSLINE1 = "general.store.addressline1"
	CONFIG_PATH_STORE_ADDRESSLINE2 = "general.store.addressline2"
	CONFIG_PATH_STORE_ZIP          = "general.store.zip"

	CONFIG_PATH_MAIL_FROM      = "general.mail.from"
	CONFIG_PATH_MAIL_SIGNATURE = "general.mail.footer"
	CONFIG_PATH_MAIL_SERVER    = "general.mail.server"
	CONFIG_PATH_MAIL_PORT      = "general.mail.port"
	CONFIG_PATH_MAIL_USER      = "general.mail.user"
	CONFIG_PATH_MAIL_PASSWORD  = "general.mail.password"
)
