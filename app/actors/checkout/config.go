package checkout

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/env"
)

func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew("can't obtain config")
	}

	// Checkout
	//---------
	err := config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Checkout",
		Description: "checkout related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	config.RegisterItem(env.T_ConfigItem{
		Path: checkout.CONFIG_PATH_CONFIRMATION_EMAIL,
		Value: `Dear {{.Visitor.last_name}} {{.Visitor.first_name}}
<br />
<br />
Thank for your order.
<br />
<h3>Order #{{.Order.increment_id}}: </h3><br />
Order summary<br />
Subtotal: ${{.Order.subtotal}}<br />
Tax: ${{.Order.tax_amount}}<br />
Shipping: ${{.Order.shipping_amount}}<br />
Total: ${{.Order.grand_total}}<br />`,
		Type:        "text",
		Editor:      "multiline_text",
		Options:     "",
		Label:       "Order confirmation e-mail: ",
		Description: "contents of email will be sent to customer on success checkout",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_CHECKOUT_TYPE,
		Value:       "accordion",
		Type:        "varchar(255)",
		Editor:      "select",
		Options:     map[string]string{"accordion": "Accordion checkout", "onepage": "OnePage checkout"},
		Label:       "Type of checkout",
		Description: "type of checkout customer will be reached by default",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Payment
	//--------
	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Payment",
		Description: "payment methods related group",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Payment Origin",
		Description: "payments methods origin information",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_COUNTRY,
		Value:       "US",
		Type:        "string",
		Editor:      "select",
		Options:     models.COUNTRIES_LIST,
		Label:       "Country",
		Description: "payment methods origin country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_STATE,
		Value:       "",
		Type:        "string",
		Editor:      "select",
		Options:     models.STATES_LIST,
		Label:       "State",
		Description: "payment methods origin state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_CITY,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "payment methods origin city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_ADDRESSLINE1,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "payment methods origin address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_ADDRESSLINE2,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "payment methods origin address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_PAYMENT_ORIGIN_ZIP,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "payment methods origin zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Shipping
	//---------
	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Shipping",
		Description: "shipping methods related group",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_GROUP,
		Value:       nil,
		Type:        env.CONFIG_ITEM_GROUP_TYPE,
		Editor:      "",
		Options:     nil,
		Label:       "Shipping Origin",
		Description: "shipping methods origin information",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_COUNTRY,
		Value:       "US",
		Type:        "string",
		Editor:      "select",
		Options:     models.COUNTRIES_LIST,
		Label:       "Country",
		Description: "shipping methods origin country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_STATE,
		Value:       "",
		Type:        "string",
		Editor:      "select",
		Options:     models.STATES_LIST,
		Label:       "State",
		Description: "shipping methods origin state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_CITY,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "shipping methods origin city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_ADDRESSLINE1,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "shipping methods origin address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_ADDRESSLINE2,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "shipping methods origin address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.T_ConfigItem{
		Path:        checkout.CONFIG_PATH_SHIPPING_ORIGIN_ZIP,
		Value:       "",
		Type:        "string",
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "shipping methods origin zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
