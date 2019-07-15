package eventbus

// GetID returns event identifier
func (it *DefaultEvent) GetID() string {
	return it.id
}

// Get retries given key from event data, if key is "" returns all data
func (it *DefaultEvent) Get(key string) interface{} {
	if key == "" {
		return it.data
	}

	if value, present := it.data[key]; present {
		return value
	}
	return nil
}

// Set updates event data with the giver key/value pair
func (it *DefaultEvent) Set(key string, value interface{}) error {
	it.data[key] = value
	return nil
}

// StopPropagation prevents event following
func (it *DefaultEvent) StopPropagation() error {
	it.propagate = false
	return nil
}
