package xmlstructs

import "encoding/xml"

// SharedStrings defines the structure of xl/sharedStrings.xml
type SharedStrings struct {
	XMLName xml.Name `xml:"http://schemas.openxmlformats.org/spreadsheetml/2006/main sst"`
	Count   int      `xml:"count,attr"`
	Unique  int      `xml:"uniqueCount,attr"`
	SI      []SI     `xml:"si"`
}

// AddString adds a string to the shared string table and returns its index.
func (ss *SharedStrings) AddString(s string) int {
	for i, si := range ss.SI {
		if si.T == s {
			return i
		}
	}
	ss.SI = append(ss.SI, SI{T: s})
	ss.Count++
	ss.Unique++
	return len(ss.SI) - 1
}
