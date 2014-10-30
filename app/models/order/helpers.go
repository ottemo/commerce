package order

import (
	"github.com/ottemo/foundation/app/models"
	"github.com/ottemo/foundation/env"
)

// retrieves current I_OrderCollection model implementation
func GetOrderCollectionModel() (I_OrderCollection, error) {
	model, err := models.GetModel(MODEL_NAME_ORDER_COLLECTION)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderCollectionModel, ok := model.(I_OrderCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_OrderCollection' capable")
	}

	return orderCollectionModel, nil
}

// retrieves current I_OrderCollection model implementation
func GetOrderItemCollectionModel() (I_OrderItemCollection, error) {
	model, err := models.GetModel(MODEL_NAME_ORDER_ITEM_COLLECTION)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderItemCollectionModel, ok := model.(I_OrderItemCollection)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_OrderItemCollection' capable")
	}

	return orderItemCollectionModel, nil
}

// retrieves current I_Order model implementation
func GetOrderModel() (I_Order, error) {
	model, err := models.GetModel(MODEL_NAME_ORDER)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	orderModel, ok := model.(I_Order)
	if !ok {
		return nil, env.ErrorNew("model " + model.GetImplementationName() + " is not 'I_Order' capable")
	}

	return orderModel, nil
}

// retrieves current I_Order model implementation and sets its ID to some value
func GetOrderModelAndSetId(orderId string) (I_Order, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderModel.SetId(orderId)
	if err != nil {
		return orderModel, env.ErrorDispatch(err)
	}

	return orderModel, nil
}

// loads order data into current I_Order model implementation
func LoadOrderById(orderId string) (I_Order, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	err = orderModel.Load(orderId)
	if err != nil {
		return nil, env.ErrorDispatch(err)
	}

	return orderModel, nil
}
