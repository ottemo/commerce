package order

// returns order item unique id, or blank string
func (it *DefaultOrderItem) GetID() string {
	return it.id
}

// sets order item unique id
func (it *DefaultOrderItem) SetID(newID string) error {
	it.id = newID
	return nil
}

// returns order item product name
func (it *DefaultOrderItem) GetName() string {
	return it.Name
}

// returns order item product sku
func (it *DefaultOrderItem) GetSku() string {
	return it.Sku
}

// returns order line item qty ordered
func (it *DefaultOrderItem) GetQty() int {
	return it.Qty
}

// returns order item product price
func (it *DefaultOrderItem) GetPrice() float64 {
	return it.Price
}

// returns order item product weight
func (it *DefaultOrderItem) GetWeight() float64 {
	return it.Weight
}

// returns order item product options
func (it *DefaultOrderItem) GetOptions() map[string]interface{} {
	return it.Options
}
