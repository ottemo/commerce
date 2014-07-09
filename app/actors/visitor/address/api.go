package address

import (
	"errors"
	"net/http"

	"github.com/ottemo/foundation/app/models"

	"github.com/ottemo/foundation/api"
)

// WEB REST API function used to obtain product attributes information
func (it *DefaultVisitorAddress) ListAddressAttributesRestAPI(resp http.ResponseWriter, req *http.Request, reqParams map[string]string, reqContent interface{}) (interface{}, error) {
	model, err := models.GetModel("VisitorAddress")
	if err != nil {
		return nil, err
	}

	address, isObject := model.(models.I_Object)
	if !isObject {
		return nil, errors.New("visitor address model is not I_Object compatible")
	}

	attrInfo := address.GetAttributesInfo()
	return attrInfo, nil
}


func (it *DefaultVisitorAddress) setupAPI() error {
	err := api.GetRestService().RegisterAPI("visitor/address", "GET", "attribute/list", it.ListAddressAttributesRestAPI)
	if err != nil {
		return err
	}

	return nil
}
