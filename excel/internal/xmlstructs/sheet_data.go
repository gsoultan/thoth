package xmlstructs

// SheetData contains the rows and cells of the worksheet
type SheetData struct {
	Rows []Row `xml:"row"`
}
