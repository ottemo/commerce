package product

import (
	"github.com/ottemo/foundation/env"
)

// Package global variables
var (
	registeredStock InterfaceStock
)

// UnRegisterStock removes stock management from system
func UnRegisterStock() error {
	registeredStock = nil
	return nil
}

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
