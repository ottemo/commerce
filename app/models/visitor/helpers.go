package visitor

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetVisitorAddressCollectionModel retrieves current InterfaceVisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (InterfaceVisitorAddressCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorAddressCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressCollectionModel, ok := model.(InterfaceVisitorAddressCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceVisitorAddressCollection' capable")
	}

	return visitorAddressCollectionModel, nil
}

// GetVisitorAddressModel retrieves current InterfaceVisitorAddress model implementation
func GetVisitorAddressModel() (InterfaceVisitorAddress, error) {
	model, err := models.GetModel(ConstModelNameVisitorAddress)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorAddressModel, ok := model.(InterfaceVisitorAddress)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceVisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// GetVisitorCollectionModel retrieves current InterfaceVisitorCollection model implementation
func GetVisitorCollectionModel() (InterfaceVisitorCollection, error) {
	model, err := models.GetModel(ConstModelNameVisitorCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorCollectionModel, ok := model.(InterfaceVisitorCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceVisitorCollection' capable")
	}

	return visitorCollectionModel, nil
}

// GetVisitorModel retrieves current InterfaceVisitor model implementation
func GetVisitorModel() (InterfaceVisitor, error) {
	model, err := models.GetModel(ConstModelNameVisitor)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	visitorModel, ok := model.(InterfaceVisitor)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceVisitor' capable")
	}

	return visitorModel, nil
}

// GetVisitorAddressModelAndSetId retrieves current InterfaceVisitorAddress model implementation and sets its ID to some value
func GetVisitorAddressModelAndSetId(visitorAddressID string) (InterfaceVisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.SetId(visitorAddressID)
	if err != nil {
		return visitorAddressModel, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// GetVisitorModelAndSetId retrieves current InterfaceVisitor model implementation and sets its ID to some value
func GetVisitorModelAndSetId(visitorID string) (InterfaceVisitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.SetId(visitorID)
	if err != nil {
		return visitorModel, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// LoadVisitorAddressById loads visitor address data into current InterfaceVisitorAddress model implementation
func LoadVisitorAddressById(visitorAddressID string) (InterfaceVisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorAddressModel.Load(visitorAddressID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorAddressModel, nil
}

// LoadVisitorById loads visitor data into current InterfaceVisitor model implementation
func LoadVisitorById(visitorID string) (InterfaceVisitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = visitorModel.Load(visitorID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return visitorModel, nil
}

// GetCurrentVisitorId returns visitor id for current session if registered or ""
func GetCurrentVisitorId(params *api.StructAPIHandlerParams) string {
	sessionVisitorID, ok := params.Session.Get(ConstSessionKeyVisitorID).(string)
	if !ok {
		return ""
	}

	return sessionVisitorID
}

// GetCurrentVisitor returns visitor for current session if registered or error
func GetCurrentVisitor(params *api.StructAPIHandlerParams) (InterfaceVisitor, error) {
	sessionVisitorID, ok := params.Session.Get(ConstSessionKeyVisitorID).(string)
	if !ok {
		return nil, env.ErrorNew("not registered visitor")
	}

	visitorInstance, err := LoadVisitorById(sessionVisitorID)

	return visitorInstance, env.ErrorDispatch(err)
}
