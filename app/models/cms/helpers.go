package cms

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// CMS Page helpers
//-----------------

// retrieves current I_CMSPage model implementation
func GetCMSPageModel() (I_CMSPage, error) {
	model, err := models.GetModel(CMS_PAGE_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	cmsPageModel, ok := model.(I_CMSPage)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_CMSPage' capable")
	}

	return cmsPageModel, nil
}

// retrieves current I_CMSPage model implementation and sets its ID to some value
func GetCMSPageModelAndSetId(cmsPageId string) (I_CMSPage, error) {

	cmsPageModel, err := GetCMSPageModel()
	if err != nil {
		return nil, err
	}

	err = cmsPageModel.SetId(cmsPageId)
	if err != nil {
		return cmsPageModel, err
	}

	return cmsPageModel, nil
}

// loads cmsPage data into current I_CMSPage model implementation
func LoadCMSPageById(cmsPageId string) (I_CMSPage, error) {

	cmsPageModel, err := GetCMSPageModel()
	if err != nil {
		return nil, err
	}

	err = cmsPageModel.Load(cmsPageId)
	if err != nil {
		return nil, err
	}

	return cmsPageModel, nil
}

// CMS Block helpers
//------------------

// retrieves current I_CMSBlock model implementation
func GetCMSBlockModel() (I_CMSBlock, error) {
	model, err := models.GetModel(CMS_BLOCK_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	csmBlockModel, ok := model.(I_CMSBlock)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_CMSBlock' capable")
	}

	return csmBlockModel, nil
}

// retrieves current I_CMSBlock model implementation and sets its ID to some value
func GetCMSBlockModelAndSetId(csmBlockId string) (I_CMSBlock, error) {

	csmBlockModel, err := GetCMSBlockModel()
	if err != nil {
		return nil, err
	}

	err = csmBlockModel.SetId(csmBlockId)
	if err != nil {
		return csmBlockModel, err
	}

	return csmBlockModel, nil
}

// loads csmBlock data into current I_CMSBlock model implementation
func LoadCMSBlockById(csmBlockId string) (I_CMSBlock, error) {

	csmBlockModel, err := GetCMSBlockModel()
	if err != nil {
		return nil, err
	}

	err = csmBlockModel.Load(csmBlockId)
	if err != nil {
		return nil, err
	}

	return csmBlockModel, nil
}
