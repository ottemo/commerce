package page

// returns page URL
func (it *DefaultCMSPage) GetURL() string {
	return it.URL
}

// returns page URL to be shown on
func (it *DefaultCMSPage) SetURL(newValue string) error {
	it.URL = newValue
	return nil
}

// returns page identifier
func (it *DefaultCMSPage) GetIdentifier() string {
	return it.Identifier
}

// sets page identifier value
func (it *DefaultCMSPage) SetIdentifier(newValue string) error {
	it.Identifier = newValue
	return nil
}

// returns page title
func (it *DefaultCMSPage) GetTitle() string {
	return it.Title
}

// sets page title value
func (it *DefaultCMSPage) SetTitle(newValue string) error {
	it.Title = newValue
	return nil
}

// returns page content
func (it *DefaultCMSPage) GetContent() string {
	return it.Content
}

// sets page content value
func (it *DefaultCMSPage) SetContent(newValue string) error {
	it.Content = newValue
	return nil
}

// returns page meta title
func (it *DefaultCMSPage) GetMetaKeywords() string {
	return it.MetaKeywords
}

// sets page meta title value
func (it *DefaultCMSPage) SetMetaKeywords(newValue string) error {
	it.MetaKeywords = newValue
	return nil
}

// returns page meta description
func (it *DefaultCMSPage) GetMetaDescription() string {
	return it.MetaDescription
}

// sets page meta description value
func (it *DefaultCMSPage) SetMetaDescription(newValue string) error {
	it.MetaDescription = newValue
	return nil
}
