package product

func (it *DefaultProduct) FromHashMap(input map[string]interface{}) error {

	if value, ok := input["_id"]; ok {
		if value, ok := value.(string); ok {
			it.id = value
		}
	}
	if value, ok := input["sku"]; ok {
		if value, ok := value.(string); ok {
			it.Sku = value
		}
	}
	if value, ok := input["name"]; ok {
		if value, ok := value.(string); ok {
			it.Name = value
		}
	}
	if value, ok := input["description"]; ok {
		if value, ok := value.(string); ok {
			it.Description = value
		}
	}
	if value, ok := input["default_image"]; ok {
		if value, ok := value.(string); ok {
			it.DefaultImage = value
		}
	}
	if value, ok := input["price"]; ok {
		if value, ok := value.(float64); ok {
			it.Price = value
		}
	}

	it.CustomAttributes.FromHashMap(input)

	return nil
}

func (it *DefaultProduct) ToHashMap() map[string]interface{} {
	result := it.CustomAttributes.ToHashMap()

	result["_id"] = it.id
	result["sku"] = it.Sku
	result["name"] = it.Name
	result["description"] = it.Name
	result["default_image"] = it.DefaultImage
	result["price"] = it.Price

	return result
}
