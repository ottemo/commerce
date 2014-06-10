package default_visitor

func (it *DefaultVisitor) FromHashMap(input map[string]interface{}) error {

	for attribute, value := range input {
		if err := it.Set(attribute, value); err != nil {
			return err
		}
	}

	return nil
}

func (it *DefaultVisitor) ToHashMap() map[string]interface{} {

	result := make( map[string]interface{} )

	result["_id"] = it.id
	result["email"] = it.Email
	result["first_name"] = it.Fname

	//result["billing_address_id"] = it.BillingAddressId
	result["billing_address"] = it.GetBillingAddress().ToHashMap()
	//result["shipping_address_id"] = it.ShippingAddressId
	result["shipping_address"] = it.GetShippingAddress().ToHashMap()

	return result
}
