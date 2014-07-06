package address

import(
	"github.com/ottemo/foundation/app/models"
)

func (it *DefaultVisitorAddress) GetModelName() string {
	return "VisitorAddress"
}

func (it *DefaultVisitorAddress) GetImplementationName() string {
	return "DefaultVisitorAddress"
}

func (it *DefaultVisitorAddress) New() (models.I_Model, error) {
	return &DefaultVisitorAddress{ }, nil
}
