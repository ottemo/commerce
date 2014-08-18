package order

import (
	"errors"
	"github.com/ottemo/foundation/app/models"
)

// retrieves current I_Order model implementation
func GetOrderModel() (I_Order, error) {
	model, err := models.GetModel(ORDER_MODEL_NAME)
	if err != nil {
		return nil, err
	}

	orderModel, ok := model.(I_Order)
	if !ok {
		return nil, errors.New("model " + model.GetImplementationName() + " is not 'I_Order' capable")
	}

	return orderModel, nil
}

// retrieves current I_Order model implementation and sets its ID to some value
func GetOrderModelAndSetId(orderId string) (I_Order, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, err
	}

	err = orderModel.SetId(orderId)
	if err != nil {
		return orderModel, err
	}

	return orderModel, nil
}

// loads order data into current I_Order model implementation
func LoadOrderById(orderId string) (I_Order, error) {

	orderModel, err := GetOrderModel()
	if err != nil {
		return nil, err
	}

	err = orderModel.Load(orderId)
	if err != nil {
		return nil, err
	}

	return orderModel, nil
}
