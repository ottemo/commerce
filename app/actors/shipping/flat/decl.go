package flat

const (
	SHIPPING_CODE = "flat_rate"
	SHIPPING_NAME = "FlatRate"

	CONFIG_PATH_GROUP = "shipping.flat_rate"

	CONFIG_PATH_ENABLED = "shipping.flat_rate.enabled"
	CONFIG_PATH_AMOUNT  = "shipping.flat_rate.amount"
	CONFIG_PATH_NAME    = "shipping.flat_rate.name"
	CONFIG_PATH_DAYS    = "shipping.flat_rate.days"
)

type FlatRateShipping struct{}
