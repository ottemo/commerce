package env

// returns config value or nil if not present
func ConfigGetValue(Path string) interface{} {
	if config := GetConfig(); config != nil {
		return config.GetValue(Path)
	}

	return nil
}
