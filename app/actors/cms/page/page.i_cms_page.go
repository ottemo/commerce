package page

// GetURL returns page URL
func (it *DefaultCMSPage) GetURL() string {
	return it.URL
}

// SetURL returns page URL to be shown on
func (it *DefaultCMSPage) SetURL(newValue string) error {
	it.URL = newValue
	return nil
}

// GetIdentifier returns page identifier
func (it *DefaultCMSPage) GetIdentifier() string {
	return it.Identifier
}

// SetIdentifier sets page identifier value
func (it *DefaultCMSPage) SetIdentifier(newValue string) error {
	it.Identifier = newValue
	return nil
}

// GetTitle returns page title
func (it *DefaultCMSPage) GetTitle() string {
	return it.Title
}

// SetTitle sets page title value
func (it *DefaultCMSPage) SetTitle(newValue string) error {
	it.Title = newValue
	return nil
}

// GetContent returns page content
func (it *DefaultCMSPage) GetContent() string {
	return it.Content
}

// SetContent sets page content value
func (it *DefaultCMSPage) SetContent(newValue string) error {
	it.Content = newValue
	return nil
}

// GetMetaKeywords returns page meta title
func (it *DefaultCMSPage) GetMetaKeywords() string {
	return it.MetaKeywords
}

// SetMetaKeywords sets page meta title value
func (it *DefaultCMSPage) SetMetaKeywords(newValue string) error {
	it.MetaKeywords = newValue
	return nil
}

// GetMetaDescription returns page meta description
func (it *DefaultCMSPage) GetMetaDescription() string {
	return it.MetaDescription
}

// SetMetaDescription sets page meta description value
func (it *DefaultCMSPage) SetMetaDescription(newValue string) error {
	it.MetaDescription = newValue
	return nil
}
