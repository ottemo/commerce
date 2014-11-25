package product

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredStock InterfaceStock
)

// RegisterStock registers given stock manager in system
func RegisterStock(stock InterfaceStock) error {
	if registeredStock != nil {
		return env.ErrorNew("stock already registered")
	}
	registeredStock = stock

	return nil
}

// GetRegisteredStock returns currently used stack manager or nil
func GetRegisteredStock() InterfaceStock {
	return registeredStock
}
