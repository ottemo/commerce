package category

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

func (it *DefaultCategory) GetName() string {
	return it.Name
}

func (it *DefaultCategory) GetProductIds() []string {
	return it.ProductIds
}

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

func (it *DefaultCategory) GetProducts() []product.InterfaceProduct {
	result := make([]product.InterfaceProduct, 0)

	for _, productId := range it.ProductIds {
		productModel, err := product.LoadProductById(productId)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}

func (it *DefaultCategory) GetParent() category.InterfaceCategory {
	return it.Parent
}

func (it *DefaultCategory) AddProduct(productId string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryId := it.GetId()
	if categoryId == "" {
		return env.ErrorNew("category ID is not set")
	}
	if productId == "" {
		return env.ErrorNew("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryId)
	collection.AddFilter("product_id", "=", productId)
	cnt, err := collection.Count()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	if cnt == 0 {
		_, err := collection.Save(map[string]interface{}{"category_id": categoryId, "product_id": productId})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	} else {
		return env.ErrorNew("junction already exists")
	}

	return nil
}

func (it *DefaultCategory) RemoveProduct(productId string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return env.ErrorNew("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(ConstCollectionNameCategoryProductJunction)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	categoryId := it.GetId()
	if categoryId == "" {
		return env.ErrorNew("category ID is not set")
	}
	if productId == "" {
		return env.ErrorNew("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryId)
	collection.AddFilter("product_id", "=", productId)
	_, err = collection.Delete()

	return env.ErrorDispatch(err)
}
