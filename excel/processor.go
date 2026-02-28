package excel

// processor handles all document manipulation operations for Excel.
// It uses composition to delegate logic to specialized smaller units.
type processor struct {
	*state
	sheetProcessor
	cellProcessor
	styleProcessor
	mediaProcessor
}
