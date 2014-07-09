package ini

import "fmt"

func (it *DefaultIniConfig) ListItems() []string {
	result := make([]string, len(it.iniFileValues))
	for itemName, _ := range it.iniFileValues {
		result = append(result, itemName)
	}
	return result
}

func (it *DefaultIniConfig) GetValue(Name string, Default string) string {

	if value, present := it.iniFileValues[Name]; present {
		return value
	} else {

		if Default == "?" {

			ok := false
			for !ok {
				fmt.Printf("%s: ", Name)
				_, err := fmt.Scanf("%s", &value)
				if err == nil {
					ok = true
				}
			}

			it.iniFileValues[Name] = value
			it.keysToStore = append(it.keysToStore, Name)

			return value
		} else {
			return Default
		}
	}
}
