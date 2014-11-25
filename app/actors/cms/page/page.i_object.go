package page

import (
	"strings"

	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"

	"github.com/ottemo/foundation/utils"
)

// Get returns object attribute value or nil
func (it *DefaultCMSPage) Get(attribute string) interface{} {
	switch strings.ToLower(attribute) {
	case "_id", "id":
		return it.GetID()
	case "url":
		return it.GetURL()
	case "identifier":
		return it.GetIdentifier()
	case "title":
		return it.GetTitle()
	case "content":
		return it.GetContent()
	case "meta_keywords":
		return it.GetMetaKeywords()
	case "meta_description":
		return it.GetMetaDescription()
	case "created_at":
		return it.CreatedAt
	case "updated_at":
		return it.UpdatedAt
	}

	return nil
}

// Set sets attribute value to object or returns error
func (it *DefaultCMSPage) Set(attribute string, value interface{}) error {
	attribute = strings.ToLower(attribute)

	switch attribute {
	case "_id", "id":
		return it.SetID(utils.InterfaceToString(value))
	case "url":
		return it.SetURL(utils.InterfaceToString(value))
	case "identifier":
		return it.SetIdentifier(utils.InterfaceToString(value))
	case "title":
		return it.SetTitle(utils.InterfaceToString(value))
	case "content":
		return it.SetContent(utils.InterfaceToString(value))
	case "meta_keywords":
		return it.SetMetaKeywords(utils.InterfaceToString(value))
	case "meta_description":
		return it.SetMetaDescription(utils.InterfaceToString(value))
	case "created_at":
		it.CreatedAt = utils.InterfaceToTime(value)
		return nil
	case "updated_at":
		it.UpdatedAt = utils.InterfaceToTime(value)
		return nil
	}

	return env.ErrorNew("unknown attribute '" + attribute + "'")
}

// FromHashMap fills object attributes from map[string]interface{}
func (it *DefaultCMSPage) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return env.ErrorDispatch(err)
		}
	}

	return nil
}

// ToHashMap represents object as map[string]interface{}
func (it *DefaultCMSPage) ToHashMap() map[string]interface{} {

	result := make(map[string]interface{})

	result["_id"] = it.GetID()
	result["url"] = it.Get("url")
	result["identifier"] = it.Get("identifier")
	result["title"] = it.Get("title")
	result["content"] = it.Get("content")
	result["meta_keywords"] = it.Get("meta_keywords")
	result["meta_description"] = it.Get("meta_description")
	result["created_at"] = it.Get("created_at")
	result["updated_at"] = it.Get("updated_at")

	return result
}

// GetAttributesInfo returns information about object attributes
func (it *DefaultCMSPage) GetAttributesInfo() []models.StructAttributeInfo {

	info := []models.StructAttributeInfo{
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "_id",
			Type:       "id",
			IsRequired: false,
			IsStatic:   true,
			Label:      "ID",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "url",
			Type:       "varchar(255)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Page URL",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "identifier",
			Type:       "varchar(255)",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Identifier",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "title",
			Type:       "varchar(255)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Title",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "content",
			Type:       "html",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Content",
			Group:      "General",
			Editors:    "html",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "meta_keywords",
			Type:       "varchar(255)",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Keywords",
			Group:      "General",
			Editors:    "line_text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "meta_description",
			Type:       "text",
			IsRequired: false,
			IsStatic:   true,
			Label:      "Meta Description",
			Group:      "General",
			Editors:    "text",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "created_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Created At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
		models.StructAttributeInfo{
			Model:      cms.ConstModelNameCMSPage,
			Collection: ConstCmsPageCollectionName,
			Attribute:  "updated_at",
			Type:       "datetime",
			IsRequired: true,
			IsStatic:   true,
			Label:      "Updated At",
			Group:      "General",
			Editors:    "not_editable",
			Options:    "",
			Default:    "",
		},
	}

	return info
}
