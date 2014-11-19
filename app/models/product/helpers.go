package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// retrieves current InterfaceProductCollection model implementation
func GetProductCollectionModel() (InterfaceProductCollection, error) {
	model, err := models.GetModel(ConstModelNameProductCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(InterfaceProductCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceProductCollection' capable")
	}

	return productModel, nil
}

// retrieves current InterfaceProduct model implementation
func GetProductModel() (InterfaceProduct, error) {
	model, err := models.GetModel(ConstModelNameProduct)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(InterfaceProduct)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceProduct' capable")
	}

	return productModel, nil
}

// retrieves current InterfaceProduct model implementation and sets its ID to some value
func GetProductModelAndSetId(productId string) (InterfaceProduct, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.SetId(productId)
	if err != nil {
		return productModel, env.ErrorDispatch(err)
	}

	return productModel, nil
}

// loads product data into current InterfaceProduct model implementation
func LoadProductById(productId string) (InterfaceProduct, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = productModel.Load(productId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return productModel, nil
}
