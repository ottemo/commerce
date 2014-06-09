package default_ini_config


func (it *DefaultIniConfig) ListItems() []string {
	result := make([]string, len(it.iniFileValues))
	for itemName, _ := range it.iniFileValues {
		result = append( result, itemName )
	}
	return result
}


func (it *DefaultIniConfig) GetValue(Name string) string {
	if value, present := it.iniFileValues[Name]; present {
		return value
	} else {
		return ""
	}
}
