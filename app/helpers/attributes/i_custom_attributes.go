package attributes

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/db"
)

// initializes helper before usage
func (it *CustomAttributes) Init(model string) (*CustomAttributes, error) {
	it.model = model
	it.values = make(map[string]interface{})

	globalCustomAttributesMutex.Lock()

	_, present := globalCustomAttributes[model]

	if present {
		it.attributes = globalCustomAttributes[model]
	} else {

		it.attributes = make(map[string]models.T_AttributeInfo)

		// retrieving information from DB
		//-------------------------------
		customAttributesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
		if err != nil {
			return it, errors.New("Can't get collection 'custom_attributes': " + err.Error())
		}

		customAttributesCollection.AddFilter("model", "=", it.model)
		dbValues, err := customAttributesCollection.Load()
		if err != nil {
			return it, errors.New("Can't load custom attributes information for '" + it.model + "'")
		}

		// filling attribute info structure
		//---------------------------------
		for _, row := range dbValues {
			attribute := models.T_AttributeInfo{}

			for key, value := range row {
				switch value := value.(type) {
				case string:
					switch key {
					case "model":
						attribute.Model = value
					case "collection":
						attribute.Collection = value
					case "attribute":
						attribute.Attribute = value
					case "type":
						attribute.Type = value
					case "label":
						attribute.Label = value
					case "group":
						attribute.Group = value
					case "editors":
						attribute.Editors = value
					case "options":
						attribute.Options = value
					case "default":
						attribute.Default = value
					case "validators":
						attribute.Validators = value
					}
				case bool:
					switch key {
					case "required":
						attribute.IsRequired = value
					case "layered":
						attribute.Layered = value
					}
				}
			}

			it.attributes[attribute.Attribute] = attribute
		}

		globalCustomAttributes[it.model] = it.attributes
	}

	globalCustomAttributesMutex.Unlock()

	return it, nil
}

// removes custom attribute from collection
func (it *CustomAttributes) RemoveAttribute(attributeName string) error {

	customAttribute, present := it.attributes[attributeName]
	if !present {
		return errors.New("There is no attribute '" + attributeName + "' for model '" + it.model + "'")
	}

	customAttributesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
	if err != nil {
		return errors.New("Can't get collection 'custom_attributes': " + err.Error())
	}

	modelCollection, err := db.GetCollection(customAttribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + customAttribute.Attribute + "' collection '" + customAttribute.Collection + "': " + err.Error())
	}

	err = modelCollection.RemoveColumn(attributeName)
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' from collection '" + customAttribute.Collection + "': " + err.Error())
	}

	globalCustomAttributesMutex.Lock()
	delete(globalCustomAttributes, it.model)
	globalCustomAttributesMutex.Unlock()

	customAttributesCollection.AddFilter("model", "=", customAttribute.Collection)
	customAttributesCollection.AddFilter("attribute", "=", attributeName)
	_, err = customAttributesCollection.Delete()
	if err != nil {
		return errors.New("Can't remove attribute '" + attributeName + "' information from 'custom_attributes' collection '" + customAttribute.Collection + "': " + err.Error())
	}

	return nil
}

// extends collection with new custom attribute
func (it *CustomAttributes) AddNewAttribute(newAttribute models.T_AttributeInfo) error {

	if _, present := it.attributes[newAttribute.Attribute]; present {
		return errors.New("There is already atribute '" + newAttribute.Attribute + "' for model '" + it.model + "'")
	}

	// getting collection where custom attribute information stores
	customAttribuesCollection, err := db.GetCollection(COLLECTION_NAME_CUSTOM_ATTRIBUTES)
	if err != nil {
		return errors.New("Can't get collection 'custom_attributes': " + err.Error())
	}

	// getting collection where attribute supposed to be
	attrCollection, err := db.GetCollection(newAttribute.Collection)
	if err != nil {
		return errors.New("Can't get attribute '" + newAttribute.Attribute + "' collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting attribute information in custom_attributes collection
	hashMap := make(map[string]interface{})

	hashMap["model"] = newAttribute.Model
	hashMap["collection"] = newAttribute.Collection
	hashMap["attribute"] = newAttribute.Attribute
	hashMap["type"] = newAttribute.Type
	hashMap["required"] = newAttribute.IsRequired
	hashMap["label"] = newAttribute.Label
	hashMap["group"] = newAttribute.Group
	hashMap["editors"] = newAttribute.Editors
	hashMap["options"] = newAttribute.Options
	hashMap["default"] = newAttribute.Default
	hashMap["validators"] = newAttribute.Validators
	hashMap["layered"] = newAttribute.Layered

	newCustomAttributeId, err := customAttribuesCollection.Save(hashMap)

	if err != nil {
		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	// inserting new attribute to supposed location
	err = attrCollection.AddColumn(newAttribute.Attribute, newAttribute.Type, false)
	if err != nil {
		customAttribuesCollection.DeleteById(newCustomAttributeId)

		return errors.New("Can't insert attribute '" + newAttribute.Attribute + "' in collection '" + newAttribute.Collection + "': " + err.Error())
	}

	it.attributes[newAttribute.Attribute] = newAttribute

	return err
}
