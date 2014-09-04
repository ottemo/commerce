package page

import (
	"errors"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/utils"
)

func setupAPI() error {

	var err error = nil

	// CMS page
	//----------

	err = api.GetRestService().RegisterAPI("cms", "GET", "page/attributes", restCMSPageAttributes)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/list", restCMSPageList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "page/list", restCMSPageList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "page/count", restCMSPageCount)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/get/:id", restCMSPageGet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "page/add", restCMSPageAdd)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "PUT", "page/update/:id", restCMSPageUpdate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "DELETE", "page/delete/:id", restCMSPageDelete)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to get CMS page available attributes information
func restCMSPageAttributes(params *api.T_APIHandlerParams) (interface{}, error) {

	cmsPage, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, err
	}

	return cmsPage.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS pages list
func restCMSPageList(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, ok := params.RequestContent.(map[string]interface{})
	if !ok {
		if params.Request.Method == "POST" {
			return nil, errors.New("unexpected request content")
		} else {
			reqData = make(map[string]interface{})
		}
	}

	// operation start
	//----------------
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handle
	cmsPageModel.ListLimit(api.GetListLimit(params))

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := cmsPageModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return cmsPageModel.List()
}

// WEB REST API function used to obtain CMS pages count in model collection
func restCMSPageCount(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(CMS_PAGE_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	return collection.Count()
}

// WEB REST API function to get CMS page information
func restCMSPageGet(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqPageId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms page id should be specified")
	}
	pageId := utils.InterfaceToString(reqPageId)

	// operation
	//----------
	cmsPage, err := cms.LoadCMSPageById(pageId)
	if err != nil {
		return nil, err
	}

	return cmsPage.ToHashMap(), nil
}

// WEB REST API for adding new CMS page in system
func restCMSPageAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetId("")
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS page in system
func restCMSPageUpdate(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	pageId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms page id should be specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// operation
	//----------
	cmsPageModel, err := cms.LoadCMSPageById(pageId)
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetId(pageId)
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS page from system
func restCMSPageDelete(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	pageId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms page id should be specified")
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModelAndSetId(pageId)
	if err != nil {
		return nil, err
	}

	cmsPageModel.Delete(pageId)

	return "ok", nil
}