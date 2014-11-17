package fedex

var (
	SHIPPING_METHODS = map[string]string{
		"EUROPE_FIRST_INTERNATIONAL_PRIORITY": "Europe First Priority",
		"FEDEX_1_DAY_FREIGHT":                 "1 Day Freight",
		"FEDEX_2_DAY_FREIGHT":                 "2 Day Freight",
		"FEDEX_2_DAY":                         "2 Day",
		"FEDEX_2_DAY_AM":                      "2 Day AM",
		"FEDEX_3_DAY_FREIGHT":                 "3 Day Freight",
		"FEDEX_EXPRESS_SAVER":                 "Express Saver",
		"FEDEX_GROUND":                        "Ground",
		"FIRST_OVERNIGHT":                     "First Overnight",
		"GROUND_HOME_DELIVERY":                "Home Delivery",
		"INTERNATIONAL_ECONOMY":               "International Economy",
		"INTERNATIONAL_ECONOMY_FREIGHT":       "Intl Economy Freight",
		"INTERNATIONAL_FIRST":                 "International First",
		"INTERNATIONAL_GROUND":                "International Ground",
		"INTERNATIONAL_PRIORITY":              "International Priority",
		"INTERNATIONAL_PRIORITY_FREIGHT":      "Intl Priority Freight",
		"PRIORITY_OVERNIGHT":                  "Priority Overnight",
		"SMART_POST":                          "Smart Post",
		"STANDARD_OVERNIGHT":                  "Standard Overnight",
		"FEDEX_FREIGHT":                       "Freight",
		"FEDEX_NATIONAL_FREIGHT":              "National Freight",
	}

	SHIPPING_DROPOFF = map[string]string{
		"REGULAR_PICKUP":          "Regular Pickup",
		"REQUEST_COURIER":         "Request Courier",
		"DROP_BOX":                "Drop Box",
		"BUSINESS_SERVICE_CENTER": "Business Service Center",
		"STATION":                 "Station",
	}

	SHIPPING_PACKAGING = map[string]string{
		"FEDEX_ENVELOPE": "FedEx Envelope",
		"FEDEX_PAK":      "FedEx Pak",
		"FEDEX_BOX":      "FedEx Box",
		"FEDEX_TUBE":     "FedEx Tube",
		"FEDEX_10KG_BOX": "FedEx 10kg Box",
		"FEDEX_25KG_BOX": "FedEx 25kg Box",
		"YOUR_PACKAGING": "Your Packaging",
	}
)

const (
	SHIPPING_CODE = "fedex"
	SHIPPING_NAME = "fedex"

	CONFIG_PATH_GROUP = "shipping.fedex"

	CONFIG_PATH_ENABLED = "shipping.fedex.enabled"
	CONFIG_PATH_TITLE   = "shipping.fedex.title"

	CONFIG_PATH_GATEWAY  = "shipping.fedex.gateway"
	CONFIG_PATH_KEY      = "shipping.fedex.key"
	CONFIG_PATH_PASSWORD = "shipping.fedex.password"
	CONFIG_PATH_NUMBER   = "shipping.fedex.number"
	CONFIG_PATH_METER    = "shipping.fedex.meter"

	CONFIG_PATH_DEFAULT_WEIGHT  = "shipping.fedex.default_weight"
	CONFIG_PATH_ALLOWED_METHODS = "shipping.fedex.allowed_methods"
	CONFIG_PATH_PACKAGING       = "shipping.fedex.packaging"
	CONFIG_PATH_DROPOFF         = "shipping.fedex.dropoff"

	CONFIG_PATH_DEBUG_LOG = "shipping.fedex.debug_log"
)

type FedEx struct{}
