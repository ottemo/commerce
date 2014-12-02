// Package stock is a default implementation of stock interface declared in
// "github.com/ottemo/foundation/app/models/checkout" package
package stock

// Package global constants
const (
	ConstCollectionNameStock = "stock"

	ConstConfigPathGroup   = "general.stock"
	ConstConfigPathEnabled = "general.stock.enabled"
)

// DefaultStock is a default implementer of InterfaceStock
type DefaultStock struct{}
