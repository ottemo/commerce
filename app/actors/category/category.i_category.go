package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

// GetName returns current category name
func (it *DefaultCategory) GetName() string {
	return it.Name
}

// GetProductIds returns product ids associated to category
func (it *DefaultCategory) GetProductIds() []string {
	return it.ProductIds
}

// GetProductsCollection returns category associated products collection instance
func (it *DefaultCategory) GetProductsCollection() product.InterfaceProductCollection {
	productCollection, err := product.GetProductCollectionModel()
	if err != nil {
		return nil
	}

	dbCollection := productCollection.GetDBCollection()
	if dbCollection != nil {
		dbCollection.AddStaticFilter("_id", "in", it.ProductIds)
	}

	return productCollection
}

// GetProducts returns a set of category associated products
func (it *DefaultCategory) GetProducts() []product.InterfaceProduct {
	var result []product.InterfaceProduct

	for _, productID := range it.ProductIds {
		productModel, err := product.LoadProductByID(productID)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}

// GetParent returns parent category of nil
func (it *DefaultCategory) GetParent() category.InterfaceCategory {
	return it.Parent
}

// AddProduct associates given product with category
func (it *DefaultCategory) AddProduct(productID string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew("category ID is not set")
	}
	if productID == "" {
		return env.ErrorNew("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryID)
	collection.AddFilter("product_id", "=", productID)
	cnt, err := collection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if cnt == 0 {
		_, err := collection.Save(map[string]interface{}{"category_id": categoryID, "product_id": productID})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew("junction already exists")
	}

	return nil
}

// RemoveProduct un-associates given product with category
func (it *DefaultCategory) RemoveProduct(productID string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryID := it.GetID()
	if categoryID == "" {
		return env.ErrorNew("category ID is not set")
	}
	if productID == "" {
		return env.ErrorNew("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryID)
	collection.AddFilter("product_id", "=", productID)
	_, err = collection.Delete()

	return env.ErrorDispatch(err)
}
