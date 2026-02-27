package xmlstructs

import "encoding/xml"

// CoreProperties defines the structure of docProps/core.xml
type CoreProperties struct {
	XMLName        xml.Name `xml:"http://schemas.openxmlformats.org/package/2006/metadata/core-properties coreProperties"`
	Title          string   `xml:"http://purl.org/dc/elements/1.1/ title,omitempty"`
	Creator        string   `xml:"http://purl.org/dc/elements/1.1/ creator,omitempty"`
	LastModifiedBy string   `xml:"http://schemas.openxmlformats.org/officeDocument/2006/metadata/core-properties lastModifiedBy,omitempty"`
	Subject        string   `xml:"http://purl.org/dc/elements/1.1/ subject,omitempty"`
	Description    string   `xml:"http://purl.org/dc/elements/1.1/ description,omitempty"`
}
