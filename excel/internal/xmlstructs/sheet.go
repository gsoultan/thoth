package xmlstructs

type Sheet struct {
	Name    string `xml:"name,attr"`
	SheetID string `xml:"sheetId,attr"`
	RID     string `xml:"r:id,attr"`
}
