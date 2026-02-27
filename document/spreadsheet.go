package document

// Spreadsheet defines operations specific to spreadsheet documents (Excel).
type Spreadsheet interface {
	Document
	// Fluent sheet selector
	Sheet(name string) (Sheet, error)

	// Additional spreadsheet-level ops
	GetSheets() ([]string, error)
}
