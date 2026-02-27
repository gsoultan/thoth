package xmlstructs

type Sheet struct {
	Name    string `xml:"name,attr"`
	SheetID string `xml:"sheetId,attr"`
	RID     string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
}
