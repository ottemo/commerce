package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"

	"github.com/ottemo/foundation/app/models/cms"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms/page/attributes", api.ConstRESTOperationGet, restCMSPageAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/list", api.ConstRESTOperationGet, restCMSPageList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/list", api.ConstRESTOperationCreate, restCMSPageList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/count", api.ConstRESTOperationGet, restCMSPageCount)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/get/:id", api.ConstRESTOperationGet, restCMSPageGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/add", api.ConstRESTOperationCreate, restCMSPageAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/update/:id", api.ConstRESTOperationUpdate, restCMSPageUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/delete/:id", api.ConstRESTOperationDelete, restCMSPageDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get CMS page available attributes information
func restCMSPageAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsPage, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsPage.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS pages list
func restCMSPageList(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	reqData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, err
	}

	// operation start
	//----------------
	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// limit parameter handle
	cmsPageCollectionModel.ListLimit(api.GetListLimit(context))

	// filters handle
	api.ApplyFilters(context, cmsPageCollectionModel.GetDBCollection())

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		cmsPageCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := cmsPageCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, env.ErrorDispatch(err)
			}
		}
	}

	return cmsPageCollectionModel.List()
}

// WEB REST API function used to obtain CMS pages count in model collection
func restCMSPageCount(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}
	dbCollection := cmsPageCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(context, dbCollection)

	// not allowing to see disabled if not admin
	if err := api.ValidateAdminRights(context); err != nil {
		dbCollection.AddFilter("enabled", "=", true)
	}

	return dbCollection.Count()
}

// WEB REST API function to get CMS page information
func restCMSPageGet(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("id")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "4c0a288a-03a4-4375-9fb3-1abd7981aecd", "cms page id should be specified")
	}

	// operation
	//----------
	cmsPage, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// not allowing to see disabled if not admin
	if api.ValidateAdminRights(context) != nil && !cmsPage.GetEnabled() {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fa76f5ac-0cce-4670-9e62-197a600ec0b9", "cms page is not available")
	}

	result := cmsPage.ToHashMap()
	result["evaluated"] = cmsPage.EvaluateContent()

	return result, nil
}

// WEB REST API for adding new CMS page in system
func restCMSPageAdd(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	reqData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID("")
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS page in system
func restCMSPageUpdate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("id")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "f128b02f-4ca5-494b-920d-5f320d112636", "cms page id should be specified")
	}

	reqData, err := api.GetRequestContentAsMap(context)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.LoadCMSPageByID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsPageModel.Set(attribute, value)
	}

	cmsPageModel.SetID(pageID)
	cmsPageModel.Save()

	return cmsPageModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS page from system
func restCMSPageDelete(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("id")
	if pageID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "20545153-b171-4638-b47e-15317e85262c", "cms page id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsPageModel, err := cms.GetCMSPageModelAndSetID(pageID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsPageModel.Delete()

	return "ok", nil
}
