package category

import (
	"errors"

	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/category"
	"github.com/ottemo/foundation/app/models/product"
)

func (it *DefaultCategory) GetName() string {
	return it.Name
}

func (it *DefaultCategory) GetProductIds() []string {
	return it.ProductIds
}

func (it *DefaultCategory) GetProductsCollection() product.I_ProductCollection {
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

func (it *DefaultCategory) GetProducts() []product.I_Product {
	result := make([]product.I_Product, 0)

	for _, productId := range it.ProductIds {
		productModel, err := product.LoadProductById(productId)
		if err == nil {
			result = append(result, productModel)
		}
	}

	return result
}

func (it *DefaultCategory) GetParent() category.I_Category {
	return it.Parent
}

func (it *DefaultCategory) AddProduct(productId string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	categoryId := it.GetId()
	if categoryId == "" {
		return errors.New("category ID is not set")
	}
	if productId == "" {
		return errors.New("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryId)
	collection.AddFilter("product_id", "=", productId)
	cnt, err := collection.Count()
	if err != nil {
		return err
	}

	if cnt == 0 {
		_, err := collection.Save(map[string]interface{}{"category_id": categoryId, "product_id": productId})
		if err != nil {
			return err
		}
	} else {
		return errors.New("junction already exists")
	}

	return nil
}

func (it *DefaultCategory) RemoveProduct(productId string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	categoryId := it.GetId()
	if categoryId == "" {
		return errors.New("category ID is not set")
	}
	if productId == "" {
		return errors.New("product ID is not set")
	}

	collection.AddFilter("category_id", "=", categoryId)
	collection.AddFilter("product_id", "=", productId)
	_, err = collection.Delete()

	return err
}
