package checkout

import(
	"github.com/ottemo/foundation/app/models/checkout"
	"github.com/ottemo/foundation/app/models"
)



// returns model name we have implementation for
func (it *DefaultCheckout) GetModelName() string {
	return checkout.CHECKOUT_MODEL_NAME
}



// returns name of current model implementation
func (it *DefaultCheckout) GetImplementationName() string {
	return "DefaultCheckout"
}


// makes new instance of model
func (it *DefaultCheckout) New() (models.I_Model, error) {
	return &DefaultCheckout{}, nil
}
