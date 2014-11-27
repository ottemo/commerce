package product

import (
	"github.com/ottemo/foundation/app/models/product"
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetID returns current product id
func (it *DefaultProduct) GetID() string {
	return it.id
}

// SetID sets current product id
func (it *DefaultProduct) SetID(NewID string) error {
	it.id = NewID
	return nil
}

// Load loads product information from DB
func (it *DefaultProduct) Load(loadID string) error {

	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecord, err := collection.LoadByID(loadID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.FromHashMap(dbRecord)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// Delete removes current product from DB
func (it *DefaultProduct) Delete() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = collection.DeleteByID(it.GetID())
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// stock management stuff
	if stockManager := product.GetRegisteredStock(); stockManager != nil {
		stockManager.RemoveProductQty(it.GetID(), make(map[string]interface{}))
	}

	return nil
}

// Save stores current product to DB
func (it *DefaultProduct) Save() error {
	collection, err := db.GetCollection(ConstCollectionNameProduct)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	valuesToStore := it.ToHashMap()
	if _, present := valuesToStore["qty"]; present {
		delete(valuesToStore, "qty")
	}
	newID, err := collection.Save(valuesToStore)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = it.SetID(newID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// stock management stuff
	if stockManager := product.GetRegisteredStock(); stockManager != nil {

		if it.qtyWasUpdated {
			err = stockManager.SetProductQty(it.GetID(), it.GetAppliedOptions(), it.Qty)
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}
		if it.optionsWereUpdated {
			err = it.saveOptionsQty()
			if err != nil {
				return env.ErrorDispatch(err)
			}
		}

	}

	return nil
}

// saveQty saves product qty to stock manager, with options qty id set
func (it *DefaultProduct) saveOptionsQty() error {

	stockManager := product.GetRegisteredStock()
	if stockManager == nil {
		return nil
	}
	//	stockManager.SetProductQty(it.GetID(), it.GetAppliedOptions(), it.Qty)

	productOptions := it.GetOptions()
	for productOptionName, productOption := range productOptions {
		if productOption, ok := productOption.(map[string]interface{}); ok {
			if qtyValue, present := productOption["qty"]; present {
				options := map[string]interface{}{productOptionName: nil}
				stockManager.SetProductQty(it.GetID(), options, utils.InterfaceToInt(qtyValue))
			}

			if productOptionValues, present := productOption["options"]; present {
				if productOptionValues, ok := productOptionValues.(map[string]interface{}); ok {

					for productOptionValueName, productOptionValue := range productOptionValues {
						if productOptionValue, ok := productOptionValue.(map[string]interface{}); ok {
							if qtyValue, present := productOptionValue["qty"]; present {
								options := map[string]interface{}{productOptionName: productOptionValueName}
								stockManager.SetProductQty(it.GetID(), options, utils.InterfaceToInt(qtyValue))
							}
						}
					}

				}
			}
		}
	}
	return nil
}
