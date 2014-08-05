package attributes

import (
	"sync"

	"github.com/ottemo/foundation/app/models"
)

const (
	CUSTOM_ATTRIBUTES_COLLECTION = "custom_attributes"
)

var (
	globalCustomAttributes      = map[string]map[string]models.T_AttributeInfo{}
	globalCustomAttributesMutex sync.RWMutex
)

type CustomAttributes struct {
	model      string
	attributes map[string]models.T_AttributeInfo

	values map[string]interface{}
}
