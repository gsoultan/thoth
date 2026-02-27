package core

import (
	"context"
	"fmt"

	"github.com/gsoultan/thoth/document"
)

// ExcelFluent provides fluent entry points for Excel documents.
type ExcelFluent struct {
	thoth *Thoth
}

func (f ExcelFluent) New(ctx context.Context) (document.Spreadsheet, error) {
	doc, err := f.thoth.createNewDocument(ctx, "temp.xlsx")
	if err != nil {
		return nil, err
	}
	return doc.(document.Spreadsheet), nil
}

func (f ExcelFluent) Open(ctx context.Context, uri string) (document.Spreadsheet, error) {
	doc, err := f.thoth.openDocument(ctx, uri)
	if err != nil {
		return nil, err
	}
	ss, ok := doc.(document.Spreadsheet)
	if !ok {
		return nil, fmt.Errorf("not an excel document: %s", uri)
	}
	return ss, nil
}

// WordFluent provides fluent entry points for Word documents.
type WordFluent struct {
	thoth *Thoth
}

func (f WordFluent) New(ctx context.Context) (document.WordProcessor, error) {
	doc, err := f.thoth.createNewDocument(ctx, "temp.docx")
	if err != nil {
		return nil, err
	}
	return doc.(document.WordProcessor), nil
}

func (f WordFluent) Open(ctx context.Context, uri string) (document.WordProcessor, error) {
	doc, err := f.thoth.openDocument(ctx, uri)
	if err != nil {
		return nil, err
	}
	wp, ok := doc.(document.WordProcessor)
	if !ok {
		return nil, fmt.Errorf("not a word document: %s", uri)
	}
	return wp, nil
}

// PDFFluent provides fluent entry points for PDF documents.
type PDFFluent struct {
	thoth *Thoth
}

func (f PDFFluent) New(ctx context.Context) (document.WordProcessor, error) {
	doc, err := f.thoth.createNewDocument(ctx, "temp.pdf")
	if err != nil {
		return nil, err
	}
	return doc.(document.WordProcessor), nil
}

func (f PDFFluent) Open(ctx context.Context, uri string) (document.WordProcessor, error) {
	doc, err := f.thoth.openDocument(ctx, uri)
	if err != nil {
		return nil, err
	}
	wp, ok := doc.(document.WordProcessor)
	if !ok {
		return nil, fmt.Errorf("not a pdf document: %s", uri)
	}
	return wp, nil
}
