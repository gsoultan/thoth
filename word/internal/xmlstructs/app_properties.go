package xmlstructs

import "encoding/xml"

// AppProperties defines the structure of docProps/app.xml
type AppProperties struct {
	XMLName     xml.Name `xml:"http://schemas.openxmlformats.org/officeDocument/2006/extended-properties Properties"`
	Application string   `xml:"Application,omitempty"`
}

// NewAppProperties creates a new instance of AppProperties with standard defaults.
func NewAppProperties() *AppProperties {
	return &AppProperties{
		Application: "Thoth Go Library",
	}
}
