package pdf

import (
	"encoding/binary"
	"fmt"
	"os"
)

type fontMetrics struct {
	name             string
	unitsPerEm       uint16
	ascent           int16
	descent          int16
	capHeight        int16
	italicAngle      float64
	isFixedPitch     bool
	numberOfHMetrics uint16
	hmtx             []byte
	cmap             []byte
	widths           map[rune]uint16
}

func parseTTF(path string) (*fontMetrics, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) < 12 {
		return nil, fmt.Errorf("invalid TTF: data too short")
	}

	numTables := binary.BigEndian.Uint16(data[4:6])
	tables := make(map[string][]byte)

	for i := range int(numTables) {
		offset := 12 + i*16
		if offset+16 > len(data) {
			break
		}
		tag := string(data[offset : offset+4])
		tableOffset := binary.BigEndian.Uint32(data[offset+8 : offset+12])
		tableLen := binary.BigEndian.Uint32(data[offset+12 : offset+16])

		if int(tableOffset+tableLen) <= len(data) {
			tables[tag] = data[tableOffset : tableOffset+tableLen]
		}
	}

	metrics := &fontMetrics{widths: make(map[rune]uint16)}

	// head table
	if head, ok := tables["head"]; ok && len(head) >= 54 {
		metrics.unitsPerEm = binary.BigEndian.Uint16(head[18:20])
	}

	// hhea table
	if hhea, ok := tables["hhea"]; ok && len(hhea) >= 36 {
		metrics.ascent = int16(binary.BigEndian.Uint16(hhea[4:6]))
		metrics.descent = int16(binary.BigEndian.Uint16(hhea[6:8]))
		metrics.numberOfHMetrics = binary.BigEndian.Uint16(hhea[34:36])
	}

	// hmtx table
	if hmtx, ok := tables["hmtx"]; ok {
		metrics.hmtx = hmtx
	}

	// cmap table
	if cmap, ok := tables["cmap"]; ok {
		metrics.cmap = cmap
		metrics.parseCmap()
	}

	// name table (simplified: just get first name)
	if name, ok := tables["name"]; ok && len(name) >= 6 {
		count := binary.BigEndian.Uint16(name[2:4])
		stringOffset := binary.BigEndian.Uint16(name[4:6])
		for i := range int(count) {
			off := 6 + i*12
			if off+12 > len(name) {
				break
			}
			nameID := binary.BigEndian.Uint16(name[off+6 : off+8])
			if nameID == 6 { // PostScript name
				length := binary.BigEndian.Uint16(name[off+8 : off+10])
				offset := binary.BigEndian.Uint16(name[off+10 : off+12])
				if int(stringOffset+offset+length) <= len(name) {
					metrics.name = string(name[stringOffset+offset : stringOffset+offset+length])
				}
				break
			}
		}
	}

	return metrics, nil
}

func (m *fontMetrics) parseCmap() {
	if len(m.cmap) < 4 {
		return
	}
	numTables := binary.BigEndian.Uint16(m.cmap[2:4])
	for i := range int(numTables) {
		off := 4 + i*8
		if off+8 > len(m.cmap) {
			break
		}
		platformID := binary.BigEndian.Uint16(m.cmap[off : off+2])
		encodingID := binary.BigEndian.Uint16(m.cmap[off+2 : off+4])
		subOffset := binary.BigEndian.Uint32(m.cmap[off+4 : off+8])

		// We prefer Unicode (Platform 0 or 3/Encoding 1)
		if (platformID == 0) || (platformID == 3 && encodingID == 1) {
			m.parseCmapTable(int(subOffset))
			return
		}
	}
}

func (m *fontMetrics) parseCmapTable(offset int) {
	if offset+6 > len(m.cmap) {
		return
	}
	format := binary.BigEndian.Uint16(m.cmap[offset : offset+2])
	if format == 4 {
		m.parseCmapFormat4(offset)
	}
}

func (m *fontMetrics) parseCmapFormat4(offset int) {
	data := m.cmap[offset:]
	if len(data) < 14 {
		return
	}
	segCount := binary.BigEndian.Uint16(data[6:8]) / 2
	endCountOff := 14
	startCountOff := endCountOff + 2 + int(segCount)*2
	idDeltaOff := startCountOff + int(segCount)*2
	idRangeOff := idDeltaOff + int(segCount)*2

	for i := range int(segCount) {
		end := binary.BigEndian.Uint16(data[endCountOff+i*2 : endCountOff+i*2+2])
		start := binary.BigEndian.Uint16(data[startCountOff+i*2 : startCountOff+i*2+2])
		delta := int16(binary.BigEndian.Uint16(data[idDeltaOff+i*2 : idDeltaOff+i*2+2]))
		rangeOffset := binary.BigEndian.Uint16(data[idRangeOff+i*2 : idRangeOff+i*2+2])

		for r := start; r <= end && r < 0xFFFF; r++ {
			var gid uint16
			if rangeOffset == 0 {
				gid = uint16(int32(r) + int32(delta))
			} else {
				// Simplified range offset handling
				addr := idRangeOff + i*2 + int(rangeOffset) + (int(r-start) * 2)
				if addr+2 <= len(data) {
					gid = binary.BigEndian.Uint16(data[addr : addr+2])
					if gid != 0 {
						gid = uint16(int32(gid) + int32(delta))
					}
				}
			}
			if gid != 0 {
				m.widths[rune(r)] = m.getGlyphWidth(gid)
			}
		}
	}
}

func (m *fontMetrics) getGlyphWidth(gid uint16) uint16 {
	if int(gid) < int(m.numberOfHMetrics) {
		off := int(gid) * 4
		if off+2 <= len(m.hmtx) {
			return binary.BigEndian.Uint16(m.hmtx[off : off+2])
		}
	} else if m.numberOfHMetrics > 0 {
		off := int(m.numberOfHMetrics-1) * 4
		if off+2 <= len(m.hmtx) {
			return binary.BigEndian.Uint16(m.hmtx[off : off+2])
		}
	}
	return 0
}
