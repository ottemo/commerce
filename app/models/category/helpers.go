package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// GetCategoryCollectionModel retrieves current InterfaceCategoryCollection model implementation
func GetCategoryCollectionModel() (InterfaceCategoryCollection, error) {
	model, err := models.GetModel(ConstModelNameCategoryCollection)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	categoryModel, ok := model.(InterfaceCategoryCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceCategoryCollection' capable")
	}

	return categoryModel, nil
}

// GetCategoryModel retrieves current InterfaceCategory model implementation
func GetCategoryModel() (InterfaceCategory, error) {
	model, err := models.GetModel(ConstModelNameCategory)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	categoryModel, ok := model.(InterfaceCategory)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'InterfaceCategory' capable")
	}

	return categoryModel, nil
}

// GetCategoryModelAndSetID retrieves current InterfaceCategory model implementation and sets its ID to some value
func GetCategoryModelAndSetID(categoryID string) (InterfaceCategory, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.SetID(categoryID)
	if err != nil {
		return categoryModel, env.ErrorDispatch(err)
	}

	return categoryModel, nil
}

// LoadCategoryByID loads category data into current InterfaceCategory model implementation
func LoadCategoryByID(categoryID string) (InterfaceCategory, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.Load(categoryID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel, nil
}
