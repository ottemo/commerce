package visitor

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_VisitorAddress model implementation
func GetVisitorAddressModel() (I_VisitorAddress, error) {
	model, err := models.GetModel(VISITOR_ADDRESS_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	visitorAddressModel, ok := model.(I_VisitorAddress)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_VisitorAddress' capable")
	}

	return visitorAddressModel, nil
}

// retrieves current I_Visitor model implementation
func GetVisitorModel() (I_Visitor, error) {
	model, err := models.GetModel(VISITOR_MODEL_NAME)
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
