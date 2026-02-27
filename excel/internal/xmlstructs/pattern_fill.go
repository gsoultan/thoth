package xmlstructs

type PatternFill struct {
	PatternType string `xml:"patternType,attr"`
	FgColor     *Color `xml:"fgColor,omitempty"`
}
