package block

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("cms/blocks", api.ConstRESTOperationGet, restCMSBlockList)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/blocks/attributes", api.ConstRESTOperationCreate, restCMSBlockAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("cms/block/:blockID", api.ConstRESTOperationGet, restCMSBlockGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/block", api.ConstRESTOperationCreate, restCMSBlockAdd)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/block/:blockID", api.ConstRESTOperationUpdate, restCMSBlockUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("cms/block/:blockID", api.ConstRESTOperationDelete, restCMSBlockDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get CMS block available attributes information
func restCMSBlockAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsBlock, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return cmsBlock.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS blocks list
//   - if "count" parameter set to non blank value returns only amount
func restCMSBlockList(context api.InterfaceApplicationContext) (interface{}, error) {

	cmsBlockCollectionModel, err := cms.GetCMSBlockCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// applying request filters
	api.ApplyFilters(context, cmsBlockCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestParameter("count") != "" {
		return cmsBlockCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	cmsBlockCollectionModel.ListLimit(api.GetListLimit(context))

	// extra parameter handle
	api.ApplyExtraAttributes(context, cmsBlockCollectionModel)

	return cmsBlockCollectionModel.List()
}

// WEB REST API function to get CMS block information
func restCMSBlockGet(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	reqBlockID := context.GetRequestArgument("blockID")
	if reqBlockID != "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a6dd2812-5070-4869-8ae2-90c4bd28bf69", "cms block id should be specified")
	}
	blockID := utils.InterfaceToString(reqBlockID)

	// operation
	//----------
	cmsBlock, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := cmsBlock.ToHashMap()
	result["evaluated"] = cmsBlock.EvaluateContent()

	return result, nil
}

// WEB REST API for adding new CMS block in system
func restCMSBlockAdd(context api.InterfaceApplicationContext) (interface{}, error) {

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
	cmsBlockModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID("")
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS block in system
func restCMSBlockUpdate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("blockID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "a7f8db95-7495-49ba-9307-baa7d5f7ecef", "cms block id should be specified")
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
	cmsBlockModel, err := cms.LoadCMSBlockByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetID(blockID)
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS block from system
func restCMSBlockDelete(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("blockID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "8dd275d4-efaf-4e67-b24d-67b28acd74e5", "cms block id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModelAndSetID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	cmsBlockModel.Delete()

	return "ok", nil
}
