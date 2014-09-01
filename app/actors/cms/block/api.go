package block

import (
	"errors"

	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/db"

	"github.com/ottemo/foundation/app/models/cms"
	"github.com/ottemo/foundation/app/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("cms", "GET", "block/attributes", restCMSBlockAttributes)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "block/list", restCMSBlockList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "block/list", restCMSBlockList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "block/count", restCMSBlockCount)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "GET", "block/get/:id", restCMSBlockGet)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "POST", "block/add", restCMSBlockAdd)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "PUT", "block/update/:id", restCMSBlockUpdate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("cms", "DELETE", "block/delete/:id", restCMSBlockDelete)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to get CMS block available attributes information
func restCMSBlockAttributes(params *api.T_APIHandlerParams) (interface{}, error) {

	cmsBlock, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, err
	}

	return cmsBlock.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain CMS pages list
func restCMSBlockList(params *api.T_APIHandlerParams) (interface{}, error) {

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
	cmsPageModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handle
	cmsPageModel.ListLimit( api.GetListLimit(params) )

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
func restCMSBlockCount(params *api.T_APIHandlerParams) (interface{}, error) {
	collection, err := db.GetCollection(CMS_BLOCK_COLLECTION_NAME)
	if err != nil {
		return nil, err
	}

	return collection.Count()
}

// WEB REST API function to get CMS block information
func restCMSBlockGet(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqBlockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms block id should be specified")
	}
	blockId := utils.InterfaceToString(reqBlockId)

	// operation
	//----------
	cmsBlock, err := cms.LoadCMSBlockById(blockId)
	if err != nil {
		return nil, err
	}

	return cmsBlock.ToHashMap(), nil
}

// WEB REST API for adding new CMS block in system
func restCMSBlockAdd(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModel()
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetId("")
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API for update existing CMS block in system
func restCMSBlockUpdate(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms block id should be specified")
	}

	reqData, err := api.GetRequestContentAsMap(params)
	if err != nil {
		return nil, err
	}

	// operation
	//----------
	cmsBlockModel, err := cms.LoadCMSBlockById( blockId )
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		cmsBlockModel.Set(attribute, value)
	}

	cmsBlockModel.SetId(blockId)
	cmsBlockModel.Save()

	return cmsBlockModel.ToHashMap(), nil
}

// WEB REST API used to delete CMS block from system
func restCMSBlockDelete(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("cms block id should be specified")
	}

	// operation
	//----------
	cmsBlockModel, err := cms.GetCMSBlockModelAndSetId(blockId)
	if err != nil {
		return nil, err
	}

	cmsBlockModel.Delete(blockId)

	return "ok", nil
}
