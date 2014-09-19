package visitor

import (
	"errors"
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_VisitorAddressCollection model implementation
func GetVisitorAddressCollectionModel() (I_VisitorAddressCollection, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_ADDRESS_COLLECTION)
	if err != nil {
		return nil, err
	}

	visitorAddressCollectionModel, ok := model.(I_VisitorAddressCollection)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_VisitorAddressCollection' capable")
	}

	return visitorAddressCollectionModel, nil
}

// retrieves current I_VisitorAddress model implementation
func GetVisitorAddressModel() (I_VisitorAddress, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_ADDRESS)
	if err != nil {
		return nil, err
	}

	visitorAddressModel, ok := model.(I_VisitorAddress)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_VisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// retrieves current I_VisitorCollection model implementation
func GetVisitorCollectionModel() (I_VisitorCollection, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR_COLLECTION)
	if err != nil {
		return nil, err
	}

	visitorCollectionModel, ok := model.(I_VisitorCollection)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_VisitorCollection' capable")
	}

	return visitorCollectionModel, nil
}

// retrieves current I_Visitor model implementation
func GetVisitorModel() (I_Visitor, error) {
	model, err := models.GetModel(MODEL_NAME_VISITOR)
	if err != nil {
		return nil, err
	}

	visitorModel, ok := model.(I_Visitor)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Visitor' capable")
	}

	return visitorModel, nil
}

// retrieves current I_VisitorAddress model implementation and sets its ID to some value
func GetVisitorAddressModelAndSetId(visitorAddressId string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, err
	}

	err = visitorAddressModel.SetId(visitorAddressId)
	if err != nil {
		return visitorAddressModel, err
	}

	return visitorAddressModel, nil
}

// retrieves current I_Visitor model implementation and sets its ID to some value
func GetVisitorModelAndSetId(visitorId string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, err
	}

	err = visitorModel.SetId(visitorId)
	if err != nil {
		return visitorModel, err
	}

	return visitorModel, nil
}

// loads visitor address data into current I_VisitorAddress model implementation
func LoadVisitorAddressById(visitorAddressId string) (I_VisitorAddress, error) {

	visitorAddressModel, err := GetVisitorAddressModel()
	if err != nil {
		return nil, err
	}

	err = visitorAddressModel.Load(visitorAddressId)
	if err != nil {
		return nil, err
	}

	return visitorAddressModel, nil
}

// loads visitor data into current I_Visitor model implementation
func LoadVisitorById(visitorId string) (I_Visitor, error) {

	visitorModel, err := GetVisitorModel()
	if err != nil {
		return nil, err
	}

	err = visitorModel.Load(visitorId)
	if err != nil {
		return nil, err
	}

	return visitorModel, nil
}

// returns visitor for current session if registered or error
func GetCurrentVisitor(params *api.T_APIHandlerParams) (I_Visitor, error) {
	sessionVisitorId, ok := params.Session.Get(SESSION_KEY_VISITOR_ID).(string)
	if !ok {
		return nil, errors.New("not registered visitor")
	}

	visitorInstance, err := LoadVisitorById(sessionVisitorId)

	return visitorInstance, err
}
