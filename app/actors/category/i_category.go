package category

import (
	"errors"

	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/app/models/category"
)

func (it *DefaultCategory) GetName() string {
	return it.Name
}


func (it *DefaultCategory) GetProducts() []product.I_Product {
	return it.Products
}


func (it *DefaultCategory) GetParent() category.I_Category {
	return it.Parent
}



func (it *DefaultCategory) AddProduct(productId string) error {

	dbEngine := db.GetDBEngine()
	if dbEngine == nil {
		return errors.New("Can't obtain DBEngine")
	}

	collection, err := dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
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
		_, err := collection.Save( map[string]interface{} { "category_id": categoryId, "product_id": productId } )
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

	collection, err := dbEngine.GetCollection(CATEGORY_PRODUCT_JUNCTION_COLLECTION_NAME)
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
