// Package attributes represents an implementation of I_CustomAttributes declared in
// "github.com/ottemo/foundation/app/models" package.
//
// In order to use it you should just embed CustomAttributes in your actor,
// you can found sample usage in "github.com/app/actors/product" package.
package attributes

import (
	"sync"

	"github.com/ottemo/foundation/app/models"
)

// Package global constants
const (
	COLLECTION_NAME_CUSTOM_ATTRIBUTES = "custom_attributes"
)

// Package global variables
var (
	globalCustomAttributes      = map[string]map[string]models.T_AttributeInfo{}
	globalCustomAttributesMutex sync.RWMutex
)

// CustomAttributes is a implementer of I_CustomAttributes
type CustomAttributes struct {
	model      string
	collection string

	attributes map[string]models.T_AttributeInfo

	values map[string]interface{}
}
