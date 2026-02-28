package pdf

import (
	"io"
	"os"

	"github.com/gsoultan/thoth/pdf/internal/objects"
)

type writeRenderer struct {
	*state
}

func (p *writeRenderer) writePDF(ctx *renderingContext, writer io.Writer, ec *objects.EncryptionContext, encryptRef *objects.Reference) error {
	// Pages
	pages := objects.Dictionary{"Type": objects.Name("Pages"), "Count": objects.Integer(len(ctx.pageRefs)), "Kids": objects.Array{}}
	pagesRef := ctx.mgr.AddObject(pages)
	for _, ref := range ctx.pageRefs {
		pages["Kids"] = append(pages["Kids"].(objects.Array), ref)
		for i := range ctx.mgr.Objects {
			if ctx.mgr.Objects[i].Number == ref.Number {
				if d, ok := ctx.mgr.Objects[i].Data.(objects.Dictionary); ok {
					d["Parent"] = pagesRef
				}
			}
		}
	}

	// Outlines (Bookmarks)
	var outlinesRef *objects.Reference
	if len(ctx.bookmarks) > 0 {
		outlinesDict := objects.Dictionary{"Type": objects.Name("Outlines"), "Count": objects.Integer(len(ctx.bookmarks))}
		outRef := ctx.mgr.AddObject(outlinesDict)
		outlinesRef = &outRef

		// Hierarchical outlines using heading levels
		type siblings struct{ first, last objects.Reference }
		children := make(map[int][]objects.Reference)    // parent.Number -> children
		prevSibling := make(map[int]objects.Reference)   // parent.Number -> last child
		parentStack := []objects.Reference{*outlinesRef} // level 0 is Outlines
		parentNumbers := []int{outlinesRef.Number}

		for _, bm := range ctx.bookmarks {
			lvl := bm.level
			if lvl < 1 {
				lvl = 1
			}
			// Ensure stack depth matches level
			for len(parentStack) > lvl {
				parentStack = parentStack[:len(parentStack)-1]
				parentNumbers = parentNumbers[:len(parentNumbers)-1]
			}
			for len(parentStack) < lvl {
				// Use last created entry as new parent
				parNum := parentNumbers[len(parentNumbers)-1]
				last, ok := prevSibling[parNum]
				if !ok {
					break
				}
				parentStack = append(parentStack, last)
				parentNumbers = append(parentNumbers, last.Number)
			}

			parent := parentStack[len(parentStack)-1]
			pageIdx := bm.page
			if pageIdx >= len(ctx.pageRefs) {
				pageIdx = len(ctx.pageRefs) - 1
			}
			if pageIdx < 0 {
				pageIdx = 0
			}
			entry := objects.Dictionary{
				"Title":  objects.PDFString(bm.title),
				"Parent": parent,
				"Dest":   objects.Array{ctx.pageRefs[pageIdx], objects.Name("XYZ"), objects.Integer(0), objects.Integer(int(bm.posY)), objects.Integer(0)},
			}
			ref := ctx.mgr.AddObject(entry)
			parNum := parent.Number
			if ps, ok := prevSibling[parNum]; ok {
				// Link Prev/Next
				for j := range ctx.mgr.Objects {
					if ctx.mgr.Objects[j].Number == ps.Number {
						d := ctx.mgr.Objects[j].Data.(objects.Dictionary)
						d["Next"] = ref
						ctx.mgr.Objects[j].Data = d
						break
					}
				}
				for j := range ctx.mgr.Objects {
					if ctx.mgr.Objects[j].Number == ref.Number {
						d := ctx.mgr.Objects[j].Data.(objects.Dictionary)
						d["Prev"] = ps
						ctx.mgr.Objects[j].Data = d
						break
					}
				}
			} else {
				children[parNum] = append(children[parNum], ref)
			}
			prevSibling[parNum] = ref
		}

		// Set First/Last for each parent having children
		for parNum, kids := range children {
			if len(kids) == 0 {
				continue
			}
			first := kids[0]
			last := kids[len(kids)-1]
			if parNum == outlinesRef.Number {
				outlinesDict["First"] = first
				outlinesDict["Last"] = last
				for j := range ctx.mgr.Objects {
					if ctx.mgr.Objects[j].Number == outRef.Number {
						ctx.mgr.Objects[j].Data = outlinesDict
						break
					}
				}
			} else {
				for j := range ctx.mgr.Objects {
					if ctx.mgr.Objects[j].Number == parNum {
						d := ctx.mgr.Objects[j].Data.(objects.Dictionary)
						d["First"] = first
						d["Last"] = last
						ctx.mgr.Objects[j].Data = d
						break
					}
				}
			}
		}
	}

	catalog := objects.Dictionary{"Type": objects.Name("Catalog"), "Pages": pagesRef}
	if outlinesRef != nil {
		catalog["Outlines"] = *outlinesRef
	}

	// Attachments (Embedded Files)
	if len(p.attachments) > 0 {
		names := objects.Dictionary{}
		embeddedFiles := objects.Dictionary{"Names": objects.Array{}}
		af := objects.Array{}

		for _, att := range p.attachments {
			data, err := os.ReadFile(att.path)
			if err != nil {
				continue
			}

			efStream := objects.Stream{
				Dict: objects.Dictionary{
					"Type":    objects.Name("EmbeddedFile"),
					"Subtype": objects.Name("application/octet-stream"), // Default
					"Params":  objects.Dictionary{"Size": objects.Integer(len(data))},
				},
				Data: data,
			}
			efRef := ctx.mgr.AddObject(efStream)

			fs := objects.Dictionary{
				"Type":           objects.Name("Filespec"),
				"F":              objects.PDFString(att.name),
				"UF":             objects.PDFString(att.name), // Unicode filename
				"EF":             objects.Dictionary{"F": efRef},
				"Desc":           objects.PDFString(att.description),
				"AFRelationship": objects.Name("Source"), // PDF 2.0
			}
			fsRef := ctx.mgr.AddObject(fs)

			embeddedFiles["Names"] = append(embeddedFiles["Names"].(objects.Array), objects.PDFString(att.name), fsRef)
			af = append(af, fsRef)
		}

		if len(af) > 0 {
			names["EmbeddedFiles"] = embeddedFiles
			catalog["Names"] = names
			catalog["AF"] = af // Associated Files (PDF 2.0)
		}
	}

	// XMP Metadata
	xmp := generateXMP(p.meta)
	metaStream := objects.Stream{
		Dict: objects.Dictionary{
			"Type":    objects.Name("Metadata"),
			"Subtype": objects.Name("XML"),
		},
		Data: []byte(xmp),
	}
	metaRef := ctx.mgr.AddObject(metaStream)
	catalog["Metadata"] = metaRef

	// Accessibility: StructTreeRoot
	hasStructs := false
	for _, pInfo := range ctx.pages {
		if len(pInfo.structItems) > 0 {
			hasStructs = true
			break
		}
	}

	if hasStructs {
		structTree := objects.Dictionary{
			"Type": objects.Name("StructTreeRoot"),
			"K":    objects.Array{},
		}
		stRef := ctx.mgr.AddObject(structTree)
		for i, pInfo := range ctx.pages {
			pageRef := ctx.pageRefs[i]
			for _, sItem := range pInfo.structItems {
				sItem["P"] = stRef
				sItem["Pg"] = pageRef
				ref := ctx.mgr.AddObject(sItem)
				structTree["K"] = append(structTree["K"].(objects.Array), ref)
			}
		}
		// Update structTree in manager
		for j := range ctx.mgr.Objects {
			if ctx.mgr.Objects[j].Number == stRef.Number {
				ctx.mgr.Objects[j].Data = structTree
				break
			}
		}
		catalog["StructTreeRoot"] = stRef
		catalog["MarkInfo"] = objects.Dictionary{"Marked": objects.Boolean(true)}
	}

	if len(ctx.allFields) > 0 {
		fields := make(objects.Array, len(ctx.allFields))
		for i, f := range ctx.allFields {
			fields[i] = f
		}
		catalog["AcroForm"] = objects.Dictionary{
			"Fields":          fields,
			"NeedAppearances": objects.Boolean(true),
			"DR": objects.Dictionary{
				"Font": objects.Dictionary{
					"Helv": objects.Dictionary{
						"Type":     objects.Name("Font"),
						"Subtype":  objects.Name("Type1"),
						"BaseFont": objects.Name("Helvetica"),
						"Encoding": objects.Name("WinAnsiEncoding"),
					},
				},
			},
			"DA": objects.PDFString("/Helv 12 Tf 0 g"),
		}
	}

	catRef := ctx.mgr.AddObject(catalog)

	var infoRef *objects.Reference
	if p.info != nil {
		ref := ctx.mgr.AddObject(p.info)
		infoRef = &ref
	}
	return ctx.mgr.Write(writer, catRef, infoRef, ec, encryptRef, ctx.fileID)
}
