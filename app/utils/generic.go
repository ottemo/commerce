package utils

// searches for presence of 1-st arg string option among provided options since 2-nd argument
func IsAmongStr(option string, searchOptions ...string) bool {
	for _, listOption := range searchOptions {
		if option == listOption {
			return true
		}
	}
	return false
}

// searches for a string in []string slice
func IsInListStr(searchItem string, searchList []string) bool {
	for _, listItem := range searchList {
		if  listItem == searchItem {
			return true
		}
	}
	return false
}

// TODO: should be somwhere in other place
func GetSiteBackUrl() string {
	return "http://ottemo:3000/"
}
