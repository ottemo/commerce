package category

import (
	"github.com/ottemo/foundation/app/utils"
	"github.com/ottemo/foundation/db"
)

func (it *DefaultCategory) GetId() string {
	return it.id
}

func (it *DefaultCategory) SetId(NewId string) error {
	it.id = NewId
	return nil
}

func (it *DefaultCategory) Load(Id string) error {

	// loading category
	categoryCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return err
	}

	dbRecord, err := categoryCollection.LoadById(Id)
	if err != nil {
		return err
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return err
	}

	it.updatePath()

	// loading category product ids
	junctionCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	junctionCollection.AddFilter("category_id", "=", it.GetId())
	junctedProducts, err := junctionCollection.Load()
	if err != nil {
		return err
	}

	for _, junctionRecord := range junctedProducts {
		it.ProductIds = append(it.ProductIds, utils.InterfaceToString(junctionRecord["product_id"]))
	}

	return nil
}

func (it *DefaultCategory) Delete() error {
	//deleting category products join
	junctionCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	err = junctionCollection.AddFilter("category_id", "=", it.GetId())
	if err != nil {
		return err
	}

	_, err = junctionCollection.Delete()
	if err != nil {
		return err
	}

	// deleting category
	categoryCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return err
	}

	err = categoryCollection.DeleteById(it.GetId())
	if err != nil {
		return err
	}

	return nil
}

func (it *DefaultCategory) Save() error {

	storingValues := it.ToHashMap()

	delete(storingValues, "products")

	categoryCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY)
	if err != nil {
		return err
	}

	// saving category
	if newId, err := categoryCollection.Save(storingValues); err == nil {
		if it.GetId() != newId {
			it.SetId(newId)
			it.updatePath()
			it.Save()
		}
	} else {
		return err
	}

	// saving category products assignment
	junctionCollection, err := db.GetCollection(COLLECTION_NAME_CATEGORY_PRODUCT_JUNCTION)
	if err != nil {
		return err
	}

	// deleting old assigned products
	junctionCollection.AddFilter("category_id", "=", it.GetId())
	_, err = junctionCollection.Delete()
	if err != nil {
		return err
	}

	// adding new assignments
	for _, categoryProductId := range it.ProductIds {
		junctionCollection.Save(map[string]interface{}{"category_id": it.GetId(), "product_id": categoryProductId})
	}

	return nil
}
