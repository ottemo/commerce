package checkout

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_Checkout model implementation
func GetCheckoutModel() (I_Checkout, error) {
	model, err := models.GetModel(CHECKOUT_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	checkoutModel, ok := model.(I_Checkout)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Checkout' capable")
	}

	return checkoutModel, nil
}
