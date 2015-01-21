package page

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms/pages", api.ConstRESTOperationGet, restCMSPageList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/pages/attributes", api.ConstRESTOperationGet, restCMSPageAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/page/:pageID", api.ConstRESTOperationGet, restCMSPageGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page", api.ConstRESTOperationCreate, restCMSPageAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/:pageID", api.ConstRESTOperationUpdate, restCMSPageUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/page/:pageID", api.ConstRESTOperationDelete, restCMSPageDelete)
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
//   - if "count" parameter set to non blank value returns only amount
func restCMSPageList(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsPageCollectionModel, err := cms.GetCMSPageCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying requested filters
	api.ApplyFilters(context, cmsPageCollectionModel.GetDBCollection())

	// excluding disabled pages for a regular visitor
	if err := api.ValidateAdminRights(context); err != nil {
		cmsPageCollectionModel.GetDBCollection().AddFilter("enabled", "=", true)
	}

	// checking for a "count" request
	if context.GetRequestParameter("count") != "" {
		return cmsPageCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	cmsPageCollectionModel.ListLimit(api.GetListLimit(context))

	// extra parameter handle
	api.ApplyExtraAttributes(context, cmsPageCollectionModel)

	return cmsPageCollectionModel.List()
}

// WEB REST API function to get CMS page information
func restCMSPageGet(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	pageID := context.GetRequestArgument("pageID")
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
	pageID := context.GetRequestArgument("pageID")
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
	pageID := context.GetRequestArgument("pageID")
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
