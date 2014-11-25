package stock

import (
	"github.com/ottemo/foundation/db"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// GetProductQty returns stock qty for a requested product-options pair
func (it *DefaultStock) GetProductQty(productID string, options map[string]interface{}) float64 {

	var minQty float64

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		env.ErrorDispatch(err)
		return minQty
	}

	// there could be couple matching request - we are looking for minimal value
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "qty", "options") {
			env.ErrorNew("unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
			continue
		}

		qty := utils.InterfaceToFloat64(dbRecord["qty"])

		if qty < minQty {
			minQty = qty
		}
	}

	return minQty
}

// RemoveProductQty removes database records matching given product-options pair
func (it *DefaultStock) RemoveProductQty(productID string, options map[string]interface{}) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// collecting record ids to remove
	var stockIDsToRemove []string
	for _, dbRecord := range dbRecords {
		if recordOptions, present := dbRecord["options"]; !present {
			recordOptions, ok := recordOptions.(map[string]interface{})

			// skipping un-matching records
			if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
				continue
			}

			if stockRecordID, present := dbRecord["_id"]; present {
				stockIDsToRemove = append(stockIDsToRemove, utils.InterfaceToString(stockRecordID))
			}
		}
	}

	// removing database records
	if len(stockIDsToRemove) > 0 {
		err := dbCollection.AddFilter("_id", "in", stockIDsToRemove)
		if err != nil {
			return env.ErrorDispatch(err)
		}

		_, err = dbCollection.Delete()
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// SetProductQty updates stock qty for a product-options pair to be exact given value
//   - use UpdateProductQty in all cases you can as it is more safer option
func (it *DefaultStock) SetProductQty(productID string, options map[string]interface{}, qty float64) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// looking for records matching request
	recordsProcessed := 0
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "_id", "qty", "options") {
			env.ErrorNew("unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
			continue
		}

		dbRecord["qty"] = qty

		_, err := dbCollection.Save(dbRecord)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		recordsProcessed++
	}

	// no records was - adding new
	if recordsProcessed == 0 {
		_, err := dbCollection.Save(map[string]interface{}{"product_id": productID, "options": options, "qty": qty})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// UpdateProductQty updates stock qty for a product-options pair on delta value which can be positive or negative
func (it *DefaultStock) UpdateProductQty(productID string, options map[string]interface{}, deltaQty float64) error {

	// receiving database information
	dbCollection, err := db.GetCollection(ConstCollectionNameStock)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = dbCollection.AddFilter("product_id", "=", productID)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	dbRecords, err := dbCollection.Load()
	if err != nil {
		return env.ErrorDispatch(err)
	}

	// looking for records matching request
	recordsProcessed := 0
	for _, dbRecord := range dbRecords {
		if !utils.StrKeysInMap(dbRecord, "_id", "qty", "options") {
			env.ErrorNew("unexpected db result")
			break
		}

		recordOptions, ok := dbRecord["options"].(map[string]interface{})

		// skipping un-matching records
		if !ok || !utils.MatchMapAValuesToMapB(recordOptions, options) {
			continue
		}

		// TODO: should work through collection update function
		qty := utils.InterfaceToFloat64(dbRecord["qty"])
		dbRecord["qty"] = qty + deltaQty

		_, err := dbCollection.Save(dbRecord)
		if err != nil {
			return env.ErrorDispatch(err)
		}
		recordsProcessed++
	}

	// no records was - adding new
	if recordsProcessed == 0 {
		_, err := dbCollection.Save(map[string]interface{}{"product_id": productID, "options": options, "qty": deltaQty})
		if err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}
