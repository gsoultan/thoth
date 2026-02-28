package xmlstructs

import (
	"encoding/xml"
	"fmt"
)

// Relationships defines the structure of OOXML .rels files
type Relationships struct {
	XMLName xml.Name       `xml:"http://schemas.openxmlformats.org/package/2006/relationships Relationships"`
	Rels    []Relationship `xml:"Relationship"`
}

// TargetByType returns the target of the first relationship of the given type
func (r *Relationships) TargetByType(relType string) string {
	for _, rel := range r.Rels {
		if rel.Type == relType {
			return rel.Target
		}
	}
	return ""
}

// AddRelationship adds a new relationship and returns its ID
func (r *Relationships) AddRelationship(relType, target string) string {
	maxID := 0
	for _, rel := range r.Rels {
		if len(rel.ID) > 3 && rel.ID[:3] == "rId" {
			var id int
			fmt.Sscanf(rel.ID, "rId%d", &id)
			if id > maxID {
				maxID = id
			}
		}
	}
	newID := fmt.Sprintf("rId%d", maxID+1)
	r.Rels = append(r.Rels, Relationship{
		ID:     newID,
		Type:   relType,
		Target: target,
	})
	return newID
}

// AddRelationshipMode adds a new relationship with TargetMode and returns its ID
func (r *Relationships) AddRelationshipMode(relType, target, mode string) string {
	maxID := 0
	for _, rel := range r.Rels {
		if len(rel.ID) > 3 && rel.ID[:3] == "rId" {
			var id int
			fmt.Sscanf(rel.ID, "rId%d", &id)
			if id > maxID {
				maxID = id
			}
		}
	}
	newID := fmt.Sprintf("rId%d", maxID+1)
	r.Rels = append(r.Rels, Relationship{
		ID:         newID,
		Type:       relType,
		Target:     target,
		TargetMode: mode,
	})
	return newID
}
