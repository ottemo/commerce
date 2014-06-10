package custom_attributes

import (
	"errors"
	"github.com/ottemo/foundation/database"
	"github.com/ottemo/foundation/models"
)

func (it *CustomAttributes) Init(model string) (*CustomAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	_, present := global_custom_attributes[model]

	if present {
		it.attributes = global_custom_attributes[model]
	} else {

		it.attributes = make(map[string]models.T_AttributeInfo)

		dbEngine := database.GetDBEngine()
		if dbEngine == nil {
			return it, errors.New("There is no database engine")
		}

		caCollection, err := dbEngine.GetCollection(CUSTOM_ATTRIBUTES_COLLECTION)
		if err != nil {
			return it, errors.New("Can't get collection 'custom_attributes': " + err.Error() )
		}

		caCollection.AddFilter("model", "=", it.model)
		dbValues, err := caCollection.Load()
		if err != nil {
			return it, errors.New("Can't load custom attributes information for '" + it.model + "'")
		}

		for _, row := range dbValues {
			attribute := models.T_AttributeInfo {
				Model:      row["model"].(string),
				Collection: row["collection"].(string),
				Attribute:  row["attribute"].(string),
				Type:       row["type"].(string),
				Label:      row["label"].(string),
				Group:      row["group"].(string),
				Editors:    row["editors"].(string),
				Options:    row["options"].(string),
				Default:    row["default"].(string),
			}

			it.attributes[attribute.Attribute] = attribute
		}

		global_custom_attributes[it.model] = it.attributes
	}

	return it, nil
}



func (it *CustomAttributes) RemoveAttribute( attributeName string ) error {

	dbEngine := database.GetDBEngine()
	if dbEngine == nil { return errors.New("There is no database engine") }

	attribute, present := it.attributes[attributeName];
	if !present {
		return errors.New("There is no attribute '" +  attributeName + "' for model '" + it.model + "'")
	}

	caCollection, err := dbEngine.GetCollection(CUSTOM_ATTRIBUTES_COLLECTION)
	if err != nil { return errors.New("Can't get collection 'custom_attributes': " + err.Error() ) }

	attrCollection, err := dbEngine.GetCollection(attribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + attribute.Attribute + "' collection '" + attribute.Collection + "': " + err.Error() )
	}

	err = attrCollection.RemoveColumn(attributeName)
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' from collection '" + attribute.Collection + "': " + err.Error() )
	}

	caCollection.AddFilter("model", "=", attribute.Collection )
	caCollection.AddFilter("attr", "=", attributeName )
	_, err = caCollection.Delete()
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' information from 'custom_attributes' collection '" + attribute.Collection + "': " + err.Error())
	}

	return nil
}



func (it *CustomAttributes) AddNewAttribute( newAttribute models.T_AttributeInfo ) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return errors.New("There is already atribute '" +  newAttribute.Attribute + "' for model '" + it.model + "'")
	}

	dbEngine := database.GetDBEngine()
	if dbEngine == nil { return errors.New("There is no database engine") }

	// getting collection where custom attribute information stores
	caCollection, err := dbEngine.GetCollection(CUSTOM_ATTRIBUTES_COLLECTION)
	if err != nil { return errors.New("Can't get collection 'custom_attributes': " + err.Error() ) }

	// getting collection where attribute supposed to be
	attrCollection, err := dbEngine.GetCollection(newAttribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + newAttribute.Attribute + "' collection '" + newAttribute.Collection + "': " + err.Error() )
	}

	// inserting attribute information in custom_attributes collection
	hashMap := make(map[string]interface{})

	hashMap["model"] = newAttribute.Model
	hashMap["collection"] = newAttribute.Collection
	hashMap["attribute"] = newAttribute.Attribute
	hashMap["type"] = newAttribute.Type
	hashMap["label"] = newAttribute.Label
	hashMap["group"] = newAttribute.Group
	hashMap["editors"] = newAttribute.Editors
	hashMap["options"] = newAttribute.Options
	hashMap["default"] = newAttribute.Default

	newCustomAttributeId, err := caCollection.Save(hashMap)

	if err != nil {
		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error() )
	}

	// inserting new attribute to supposed location
	err = attrCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		caCollection.DeleteById(newCustomAttributeId)

		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error() )
	}

	it.attributes[newAttribute.Attribute] = newAttribute

	return err
}
