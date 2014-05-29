package model_product

import ("errors")

// Interfaces declaration
//-----------------------

type I_ProductModel interface {
	GetModelName() string
}


// Delegate routines
//------------------
var registeredProductModel I_WebServer

func GetModel() I_ProductModel {
	return registeredProductModel
}

func RegisterModel(ProductModel I_ProductModel) error {
	if registeredProductModel == nil {
		registeredProductModel = ProductModel
	} else {
		return errors.New("Product model '" + registeredProductModel.GetName() + "' already registered")
	}
	return nil
}
