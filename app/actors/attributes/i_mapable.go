package attributes

func (it *CustomAttributes) FromHashMap(input map[string]interface{}) error {
	it.values = input
	return nil
}

func (it *CustomAttributes) ToHashMap() map[string]interface{} {
	return it.values
}
