package product

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// retrieves current I_ProductCollection model implementation
func GetProductCollectionModel() (I_ProductCollection, error) {
	model, err := models.GetModel(MODEL_NAME_PRODUCT_COLLECTION)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(I_ProductCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_ProductCollection' capable")
	}

	return productModel, nil
}

// retrieves current I_Product model implementation
func GetProductModel() (I_Product, error) {
	model, err := models.GetModel(MODEL_NAME_PRODUCT)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	productModel, ok := model.(I_Product)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Product' capable")
	}

	return productModel, nil
}

// retrieves current I_Product model implementation and sets its ID to some value
func GetProductModelAndSetId(productId string) (I_Product, error) {

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

// loads product data into current I_Product model implementation
func LoadProductById(productId string) (I_Product, error) {

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
