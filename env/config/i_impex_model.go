package config

import (
	"github.com/ottemo/foundation/env"
	"github.com/ottemo/foundation/utils"
)

// Import imports value in config through Impex
func (it *DefaultConfig) Import(item map[string]interface{}, testMode bool) (map[string]interface{}, error) {
	var key string
	var value interface{}

	if pathValue, present := item["path"]; present {
		key = utils.InterfaceToString(pathValue)
	} else if keyValue, present := item["key"]; present {
		key = utils.InterfaceToString(keyValue)
	}

	value, _ = item["value"]

	if testMode == false {
		err := it.SetValue(key, value)
		return item, err
	}

	return item, nil
}

// Export exports config values through Impex
func (it *DefaultConfig) Export(iterator func(map[string]interface{}) bool) error {
	for _, itemPath := range it.ListPathes() {
		itemInfo := it.GetItemsInfo(itemPath)
		if len(itemInfo) > 0 {
			if itemInfo[0].Type == env.ConstConfigItemGroupType {
				continue
			}
		}
		value := it.GetValue(itemPath)

		continueFlag := iterator(map[string]interface{}{"path": itemPath, "value": value})
		if continueFlag == false {
			break
		}
	}

	return nil
}
