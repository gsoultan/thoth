package xmlstructs

import "encoding/xml"

// ContentTypes defines the structure of [Content_Types].xml
type ContentTypes struct {
	XMLName  xml.Name   `xml:"http://schemas.openxmlformats.org/package/2006/content-types Types"`
	Defaults []Default  `xml:"Default"`
	Override []Override `xml:"Override"`
}

// Default defines a default content type for a specific extension
type Default struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// Override defines a specific content type for a specific part
type Override struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

// NewContentTypes creates a new instance of ContentTypes with standard defaults.
func NewContentTypes() *ContentTypes {
	return &ContentTypes{
		Defaults: []Default{
			{Extension: "rels", ContentType: "application/vnd.openxmlformats-package.relationships+xml"},
			{Extension: "xml", ContentType: "application/xml"},
		},
		Override: []Override{
			{PartName: "/xl/workbook.xml", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"},
			{PartName: "/xl/styles.xml", ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"},
		},
	}
}

// AddOverride adds or updates an override for a part.
func (ct *ContentTypes) AddOverride(partName, contentType string) {
	for i, o := range ct.Override {
		if o.PartName == partName {
			ct.Override[i].ContentType = contentType
			return
		}
	}
	ct.Override = append(ct.Override, Override{PartName: partName, ContentType: contentType})
}

// AddDefault adds a default content type for an extension if it doesn't already exist.
func (ct *ContentTypes) AddDefault(extension, contentType string) {
	for i, d := range ct.Defaults {
		if d.Extension == extension {
			ct.Defaults[i].ContentType = contentType
			return
		}
	}
	ct.Defaults = append(ct.Defaults, Default{Extension: extension, ContentType: contentType})
}
