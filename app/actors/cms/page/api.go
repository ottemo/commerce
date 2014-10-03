package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cms"
)

func setupAPI() error {

	var err error = nil

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
	err = api.GetRestService().RegisterAPI("cms", "GET", "page/count", restCMSPageCount)
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
			return nil, env.ErrorNew("unexpected request content")
		} else {
			reqData = make(map[string]interface{})
		}
	}

	// operation start
	//----------------
	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handle
	cmsPageCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, cmsPageCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := cmsPageCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return cmsPageCollectionModel.List()
}

// WEB REST API function used to obtain CMS pages count in model collection
func restCMSPageCount(params *api.T_APIHandlerParams) (interface{}, error) {

	cmsPageCollectionModel, err := cms.GetCMSBlockCollectionModel()
	if err != nil {
		return nil, err
	}
	dbCollection := cmsPageCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function to get CMS page information
func restCMSPageGet(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqPageId, present := params.RequestURLParams["id"]
	if !present {
		return nil, env.ErrorNew("cms page id should be specified")
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

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
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
		return nil, env.ErrorNew("cms page id should be specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
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
		return nil, env.ErrorNew("cms page id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModelAndSetId(pageId)
	if err != nil {
		return nil, err
	}

	cmsPageModel.Delete()

	return "ok", nil
}
