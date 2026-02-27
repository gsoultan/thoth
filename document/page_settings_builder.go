package document

// PageSettingsBuilder helps in constructing PageSettings using a chainable API.
type PageSettingsBuilder struct {
	settings PageSettings
}

// NewPageSettingsBuilder creates a new PageSettingsBuilder.
func NewPageSettingsBuilder() *PageSettingsBuilder {
	return &PageSettingsBuilder{}
}

// WithOrientation sets the orientation for the page settings.
func (b *PageSettingsBuilder) WithOrientation(o Orientation) *PageSettingsBuilder {
	b.settings.Orientation = o
	return b
}

// WithPaperType sets the paper type for the page settings.
func (b *PageSettingsBuilder) WithPaperType(p PaperType) *PageSettingsBuilder {
	b.settings.PaperType = p
	return b
}

// WithMargins sets the margins for the page settings.
func (b *PageSettingsBuilder) WithMargins(m Margins) *PageSettingsBuilder {
	b.settings.Margins = m
	return b
}

// Build returns the constructed PageSettings.
func (b *PageSettingsBuilder) Build() PageSettings {
	return b.settings
}
