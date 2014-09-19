package order

import (
	"errors"

	"github.com/ottemo/foundation/api"

	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/utils"
)

func setupAPI() error {

	var err error = nil

	err = api.GetRestService().RegisterAPI("order", "GET", "attributes", restOrderAttributes)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "list", restOrderList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("order", "POST", "list", restOrderList)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "count", restOrderCount)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("order", "GET", "get/:id", restOrderGet)
	if err != nil {
		return err
	}
	// err = api.GetRestService().RegisterAPI("order", "POST", "add", restOrderAdd)
	// if err != nil {
	// 	return err
	// }
	err = api.GetRestService().RegisterAPI("order", "PUT", "update/:id", restOrderUpdate)
	if err != nil {
		return err
	}
	err = api.GetRestService().RegisterAPI("order", "DELETE", "delete/:id", restOrderDelete)
	if err != nil {
		return err
	}

	return nil
}

// WEB REST API function to get order available attributes information
func restOrderAttributes(params *api.T_APIHandlerParams) (interface{}, error) {

	orderModel, err := order.GetOrderModel()
	if err != nil {
		return nil, err
	}

	return orderModel.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain orders list
func restOrderList(params *api.T_APIHandlerParams) (interface{}, error) {

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

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	// operation start
	//----------------
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, err
	}

	// limit parameter handle
	orderCollectionModel.ListLimit(api.GetListLimit(params))

	// filters handle
	api.ApplyFilters(params, orderCollectionModel.GetDBCollection())

	// extra parameter handle
	if extra, isExtra := reqData["extra"]; isExtra {
		extra := utils.Explode(utils.InterfaceToString(extra), ",")
		for _, value := range extra {
			err := orderCollectionModel.ListAddExtraAttribute(value)
			if err != nil {
				return nil, err
			}
		}
	}

	return orderCollectionModel.List()
}

// WEB REST API function used to obtain orders count in model collection
func restOrderCount(params *api.T_APIHandlerParams) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, err
	}
	dbCollection := orderCollectionModel.GetDBCollection()

	// filters handle
	api.ApplyFilters(params, dbCollection)

	return dbCollection.Count()
}

// WEB REST API function to get order information
func restOrderGet(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	reqBlockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("order id should be specified")
	}
	blockId := utils.InterfaceToString(reqBlockId)

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	// operation
	//----------
	orderModel, err := order.LoadOrderById(blockId)
	if err != nil {
		return nil, err
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()
	return result, nil
}

// WEB REST API for update existing order in system
func restOrderUpdate(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("order id should be specified")
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
	orderModel, err := order.LoadOrderById(blockId)
	if err != nil {
		return nil, err
	}

	for attribute, value := range reqData {
		orderModel.Set(attribute, value)
	}

	orderModel.SetId(blockId)
	orderModel.Save()

	return orderModel.ToHashMap(), nil
}

// WEB REST API used to delete order from system
func restOrderDelete(params *api.T_APIHandlerParams) (interface{}, error) {

	// check request params
	//---------------------
	blockId, present := params.RequestURLParams["id"]
	if !present {
		return nil, errors.New("order id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(params); err != nil {
		return nil, err
	}

	// operation
	//----------
	orderModel, err := order.GetOrderModelAndSetId(blockId)
	if err != nil {
		return nil, err
	}

	orderModel.Delete()

	return "ok", nil
}
