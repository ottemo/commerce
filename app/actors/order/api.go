package order

import (
	"github.com/ottemo/foundation/api"
	"github.com/ottemo/foundation/app/models/order"
	"github.com/ottemo/foundation/env"
)

// setupAPI setups package related API endpoint routines
func setupAPI() error {

	var err error

	err = api.GetRestService().RegisterAPI("orders/attributes", api.ConstRESTOperationGet, restOrderAttributes)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("orders", api.ConstRESTOperationGet, restOrderList)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	err = api.GetRestService().RegisterAPI("order/:orderID", api.ConstRESTOperationGet, restOrderGet)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	// err = api.GetRestService().RegisterAPI("order", api.ConstRESTOperationCreate, restOrderAdd)
	// if err != nil {
	// 	return env.ErrorDispatch(err)
	// }
	err = api.GetRestService().RegisterAPI("order/:orderID", api.ConstRESTOperationUpdate, restOrderUpdate)
	if err != nil {
		return env.ErrorDispatch(err)
	}
	err = api.GetRestService().RegisterAPI("order/:orderID", api.ConstRESTOperationDelete, restOrderDelete)
	if err != nil {
		return env.ErrorDispatch(err)
	}

	return nil
}

// WEB REST API function to get order available attributes information
func restOrderAttributes(context api.InterfaceApplicationContext) (interface{}, error) {

	orderModel, err := order.GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel.GetAttributesInfo(), nil
}

// WEB REST API function used to obtain orders list
//   - if "count" parameter set to non blank value returns only amount
func restOrderList(context api.InterfaceApplicationContext) (interface{}, error) {

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// taking orders collection model
	orderCollectionModel, err := order.GetOrderCollectionModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// filters handle
	api.ApplyFilters(context, orderCollectionModel.GetDBCollection())

	// checking for a "count" request
	if context.GetRequestParameter("count") != "" {
		return orderCollectionModel.GetDBCollection().Count()
	}

	// limit parameter handle
	orderCollectionModel.ListLimit(api.GetListLimit(context))

	// extra parameter handle
	api.ApplyExtraAttributes(context, orderCollectionModel)

	return orderCollectionModel.List()
}

// WEB REST API function to get order information
func restOrderGet(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("orderID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "723ef443-f974-4455-9be0-a8af13916554", "order id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.LoadOrderByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	result := orderModel.ToHashMap()
	result["items"] = orderModel.GetItems()
	return result, nil
}

// WEB REST API for update existing order in system
func restOrderUpdate(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("orderID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "20a08638-e9e6-428b-b70c-a418d7821e4b", "order id should be specified")
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
	orderModel, err := order.LoadOrderByID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	for attribute, value := range reqData {
		orderModel.Set(attribute, value)
	}

	orderModel.SetID(blockID)
	orderModel.Save()

	return orderModel.ToHashMap(), nil
}

// WEB REST API used to delete order from system
func restOrderDelete(context api.InterfaceApplicationContext) (interface{}, error) {

	// check request context
	//---------------------
	blockID := context.GetRequestArgument("orderID")
	if blockID == "" {
		return nil, env.ErrorNew(ConstErrorModule, env.ConstErrorLevelAPI, "fc3011c7-e58c-4433-b9b0-881a7ba005cf", "order id should be specified")
	}

	// check rights
	if err := api.ValidateAdminRights(context); err != nil {
		return nil, env.ErrorDispatch(err)
	}

	// operation
	//----------
	orderModel, err := order.GetOrderModelAndSetID(blockID)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel.Delete()

	return "ok", nil
}
