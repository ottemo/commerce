package visitor

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

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
