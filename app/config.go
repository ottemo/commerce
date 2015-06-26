package app

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupConfig setups package configuration values for a system
func setupConfig() error {
	config := env.GetConfig()
	if config == nil {
		return env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "f635e96c-3cd7-4ae2-a507-4349021e9f13", "can't obtain config")
	}

	err := config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "General",
		Description: "application general options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathAppGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Application",
		Description: "application related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStorefrontURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Storefront host URL",
		Description: "URL application will use to generate storefront resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathDashboardURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Dashboard host URL",
		Description: "URL application will use to generate dashboard resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathFoundationURL,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Foundation host URL",
		Description: "URL application will use to generate foundation resources links",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Store",
		Description: "web store related options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreName,
		Value:       "Ottemo store",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Name",
		Description: "name of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreEmail,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "E-mail",
		Description: "e-mail of your web store",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreRootLogin,
		Value:       "admin",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Root login",
		Description: "login to enter admin panel as root",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreRootPassword,
		Value:       "admin",
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     nil,
		Label:       "Root password",
		Description: "password to enter admin panel as root",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreCountry,
		Value:       "US",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstCountriesList,
		Label:       "Country",
		Description: "store location country",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreState,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "select",
		Options:     models.ConstStatesList,
		Label:       "State",
		Description: "store location state",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreCity,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "City",
		Description: "store location city",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreAddressline1,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 1",
		Description: "store location address line 1",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreAddressline2,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "Address Line 2",
		Description: "store location address line 2",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathStoreZip,
		Value:       "",
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "line_text",
		Options:     "",
		Label:       "zip",
		Description: "store location zip code",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailGroup,
		Value:       nil,
		Type:        env.ConstConfigTypeGroup,
		Editor:      "",
		Options:     nil,
		Label:       "Mail",
		Description: "web store mailing options",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailServer,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "Host",
		Description: "web store mailing server",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailPort,
		Value:       587,
		Type:        env.ConstConfigTypeInteger,
		Editor:      "integer",
		Options:     nil,
		Label:       "Port",
		Description: "web store mailing server port",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailUser,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "User",
		Description: "mail server username",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailPassword,
		Value:       nil,
		Type:        env.ConstConfigTypeSecret,
		Editor:      "password",
		Options:     nil,
		Label:       "Password",
		Description: "mail server password",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailFrom,
		Value:       nil,
		Type:        env.ConstConfigTypeVarchar,
		Editor:      "text",
		Options:     nil,
		Label:       "From",
		Description: "sending mail from field",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathMailSignature,
		Value:       "Sincerely, Ottemo",
		Type:        env.ConstConfigTypeText,
		Editor:      "multiline_text",
		Options:     nil,
		Label:       "Signature",
		Description: "sending mail signature",
		Image:       "",
	}, nil)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	// Hide levels
	emailValidator := func(newValue interface{}) (interface{}, error) {
		newEmail := utils.InterfaceToString(newValue)
		if newEmail == "" {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "d5abe68b-5bde-4b14-a3a7-b89507c14597", "recipient e-mail can not be blank")
			return "support+ContactUs@ottemo.io", err
		}
		if !utils.ValidEmailAddress(newEmail) {
			err := env.ErrorNew(ConstErrorModule, env.ConstErrorLevelStartStop, "720c5b82-4aa9-405b-b52d-b94d1f31e49d", "recipient e-mail is not in a valid format")
			return "support+ContactUs@ottemo.io", err

		}

		return newEmail, nil
	}
	err = config.RegisterItem(env.StructConfigItem{
		Path:        ConstConfigPathContactUsRecipient,
		Value:       "support+ContactUs@ottemo.io",
		Type:        utils.DataTypeWPrecision(env.ConstConfigTypeVarchar, 255),
		Editor:      "text",
		Options:     nil,
		Label:       "Contact Us Recipient",
		Description: "email of contact form recipient",
		Image:       "",
	}, emailValidator)

	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}
