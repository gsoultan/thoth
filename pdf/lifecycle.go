package pdf

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/gsoultan/thoth/pdf/internal/objects"
	"github.com/gsoultan/thoth/pdf/internal/parser"
)

// lifecycle handles document lifecycle operations.
type lifecycle struct{ *state }

// Open loads a document from a reader.
func (p *lifecycle) Open(ctx context.Context, reader io.Reader) error {
	l := parser.NewLexer(reader)
	pr := parser.NewParser(l)

	for {
		obj, err := pr.ParseObject()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			// Skip single non-objects (e.g., trailer keyword, %%EOF)
			continue
		}
		if obj == nil {
			// Check if we hit EOF via nil obj
			break
		}
		p.objects = append(p.objects, obj)
	}

	return nil
}

// Save writes the document to a writer.
func (p *lifecycle) Save(ctx context.Context, writer io.Writer) error {
	if len(p.objects) > 0 {
		return p.saveModified(ctx, writer)
	}

	hasTOC := false
	for _, item := range p.contentItems {
		if item.isTOC {
			hasTOC = true
			break
		}
	}

	r := &renderer{p.state}
	pr := &pageRenderer{p.state}
	wr := &writeRenderer{p.state}

	renderCtx := r.newRenderingContext()

	if hasTOC {
		// Pass 1: Dry run to collect bookmarks/headings
		for name, path := range p.fonts {
			r.ensureFontInContext(renderCtx, name, path)
		}
		r.collectImages(renderCtx, p.contentItems)
		r.collectImages(renderCtx, p.header)
		r.collectImages(renderCtx, p.footer)

		pr.renderContent(renderCtx)
		pr.finishPage(renderCtx)

		// Reset context for real run, but KEEP collected bookmarks
		bookmarks := renderCtx.bookmarks
		renderCtx = r.newRenderingContext()
		renderCtx.bookmarks = bookmarks
	}

	// Real run
	for name, path := range p.fonts {
		r.ensureFontInContext(renderCtx, name, path)
	}

	// Pre-process images
	r.collectImages(renderCtx, p.contentItems)
	r.collectImages(renderCtx, p.header)
	r.collectImages(renderCtx, p.footer)

	// Render content
	pr.renderContent(renderCtx)

	pr.finishPage(renderCtx)

	pr.finalizePages(renderCtx)

	var ec *objects.EncryptionContext
	var encryptRef *objects.Reference
	fileID := make([]byte, 16)
	_, _ = rand.Read(fileID)

	if p.password != "" {
		ec = objects.NewEncryptionContext(p.password, fileID)

		encryptDict := objects.Dictionary{
			"Filter":          objects.Name("Standard"),
			"V":               objects.Integer(ec.Algorithm),
			"R":               objects.Integer(ec.Revision),
			"O":               objects.PDFString(ec.O),
			"U":               objects.PDFString(ec.U),
			"P":               objects.Integer(int(ec.P)),
			"EncryptMetadata": objects.Boolean(true),
		}

		if ec.Algorithm >= 4 {
			// AES requires CryptFilter
			cfName := objects.Name("StdCF")
			cfm := objects.Name("AESV2")
			if ec.Algorithm == 5 {
				cfm = objects.Name("AESV3")
			}
			encryptDict["CF"] = objects.Dictionary{
				"StdCF": objects.Dictionary{
					"Type":   objects.Name("CryptFilter"),
					"CFM":    cfm,
					"Length": objects.Integer(len(ec.EncryptKey) * 8),
				},
			}
			encryptDict["StmF"] = cfName
			encryptDict["StrF"] = cfName
		} else {
			encryptDict["Length"] = objects.Integer(len(ec.EncryptKey) * 8)
		}

		ref := renderCtx.mgr.AddObject(encryptDict)
		encryptRef = &ref
	}
	renderCtx.fileID = fileID

	if err := wr.writePDF(renderCtx, writer, ec, encryptRef); err != nil {
		return fmt.Errorf("write pdf: %w", err)
	}
	return nil
}

// Close releases any resources used by the document.
func (p *lifecycle) Close() error {
	return nil
}
