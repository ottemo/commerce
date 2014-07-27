package product

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_Product model implementation
func GetProductModel() (I_Product, error) {
	model, err := models.GetModel(PRODUCT_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	productModel, ok := model.(I_Product)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Product' capable")
	}

	return productModel, nil
}


// retrieves current I_Product model implementation and sets its ID to some value
func GetProductModelAndSetId(productId string) (I_Product, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, err
	}

	err = productModel.SetId(productId)
	if err != nil {
		return productModel, err
	}

	return productModel, nil
}



// loads product data into current I_Product model implementation
func LoadProductById(productId string) (I_Product, error) {

	productModel, err := GetProductModel()
	if err != nil {
		return nil, err
	}

	err = productModel.Load(productId)
	if err != nil {
		return nil, err
	}

	return productModel, nil
}
