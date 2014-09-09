package category

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_CategoryCollection model implementation
func GetCategoryCollectionModel() (I_CategoryCollection, error) {
	model, err := models.GetModel(MODEL_NAME_CATEGORY_COLLECTION)
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(I_CategoryCollection)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_CategoryCollection' capable")
	}

	return categoryModel, nil
}

// retrieves current I_Category model implementation
func GetCategoryModel() (I_Category, error) {
	model, err := models.GetModel(MODEL_NAME_CATEGORY)
	if err != nil {
		return nil, err
	}

	categoryModel, ok := model.(I_Category)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Category' capable")
	}

	return categoryModel, nil
}

// retrieves current I_Category model implementation and sets its ID to some value
func GetCategoryModelAndSetId(categoryId string) (I_Category, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, err
	}

	err = categoryModel.SetId(categoryId)
	if err != nil {
		return categoryModel, err
	}

	return categoryModel, nil
}

// loads category data into current I_Category model implementation
func LoadCategoryById(categoryId string) (I_Category, error) {

	categoryModel, err := GetCategoryModel()
	if err != nil {
		return nil, err
	}

	err = categoryModel.Load(categoryId)
	if err != nil {
		return nil, err
	}

	return categoryModel, nil
}
