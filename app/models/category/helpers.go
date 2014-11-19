package category

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// retrieves current InterfaceCategoryCollection model implementation
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

// retrieves current InterfaceCategory model implementation
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

// retrieves current InterfaceCategory model implementation and sets its ID to some value
func GetCategoryModelAndSetId(categoryId string) (InterfaceCategory, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.SetId(categoryId)
	if err != nil {
		return categoryModel, env.ErrorDispatch(err)
	}

	return categoryModel, nil
}

// loads category data into current InterfaceCategory model implementation
func LoadCategoryById(categoryId string) (InterfaceCategory, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = categoryModel.Load(categoryId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return categoryModel, nil
}
