package excel

import (
	"archive/zip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/gsoultan/thoth/excel/internal/xmlstructs"
)

func (e *state) loadCore(ctx context.Context) error {
	// 0. Load Content Types
	var ct xmlstructs.ContentTypes
	if err := e.loadXML("[Content_Types].xml", &ct); err == nil {
		e.contentTypes = &ct
	} else {
		e.contentTypes = xmlstructs.NewContentTypes()
	}

	// 1. Load root relationships to find workbook
	var rootRels xmlstructs.Relationships
	if err := e.loadXML("_rels/.rels", &rootRels); err != nil {
		return fmt.Errorf("load root rels: %w", err)
	}
	e.rootRels = &rootRels
	e.rootRelsPath = "_rels/.rels"

	workbookPath := rootRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument")
	if workbookPath == "" {
		workbookPath = "xl/workbook.xml" // fallback
	} else if strings.HasPrefix(workbookPath, "/") {
		workbookPath = workbookPath[1:]
	}

	// 2. Load workbook
	var wb xmlstructs.Workbook
	if err := e.loadXML(workbookPath, &wb); err != nil {
		return fmt.Errorf("load workbook: %w", err)
	}
	e.workbook = &wb

	// 3. Load workbook relationships to find other parts
	wbRelsPath := strings.Replace(workbookPath, "workbook.xml", "_rels/workbook.xml.rels", 1)
	e.wbRelsPath = wbRelsPath
	var wbRels xmlstructs.Relationships
	e.loadXML(wbRelsPath, &wbRels)
	e.workbookRels = &wbRels

	// Shared Strings
	ssPath := wbRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings")
	if ssPath != "" {
		if !strings.HasPrefix(ssPath, "/") {
			ssPath = "xl/" + ssPath // simplified base path handling
		} else {
			ssPath = ssPath[1:]
		}
		var ss xmlstructs.SharedStrings
		if err := e.loadXML(ssPath, &ss); err == nil {
			e.sharedStrings = &ss
		}
	}

	// Core Properties
	cpPath := rootRels.TargetByType("http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties")
	if cpPath != "" {
		if strings.HasPrefix(cpPath, "/") {
			cpPath = cpPath[1:]
		}
		var cp xmlstructs.CoreProperties
		if err := e.loadXML(cpPath, &cp); err == nil {
			e.coreProperties = &cp
		}
	}

	// Styles
	stylesPath := wbRels.TargetByType("http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles")
	if stylesPath != "" {
		if !strings.HasPrefix(stylesPath, "/") {
			stylesPath = "xl/" + stylesPath
		} else {
			stylesPath = stylesPath[1:]
		}
		var s xmlstructs.Styles
		if err := e.loadXML(stylesPath, &s); err == nil {
			e.styles = &s
		}
	}

	// Load sheets
	for _, s := range e.workbook.Sheets {
		target := ""
		for _, rel := range e.workbookRels.Rels {
			if rel.ID == s.RID {
				target = rel.Target
				break
			}
		}

		if target == "" {
			target = fmt.Sprintf("worksheets/sheet%s.xml", s.SheetID)
		}

		var path string
		if strings.HasPrefix(target, "/") {
			path = target[1:]
		} else {
			path = "xl/" + target
		}

		var ws xmlstructs.Worksheet
		if err := e.loadXML(path, &ws); err != nil {
			continue
		}
		e.sheets[s.Name] = &ws

		// Load sheet rels
		relPath := ""
		if strings.Contains(path, "/") {
			lastSlash := strings.LastIndex(path, "/")
			relPath = path[:lastSlash] + "/_rels/" + path[lastSlash+1:] + ".rels"
		} else {
			relPath = "_rels/" + path + ".rels"
		}

		var wRels xmlstructs.Relationships
		if err := e.loadXML(relPath, &wRels); err == nil {
			e.sheetRels[s.Name] = &wRels
		} else {
			e.sheetRels[s.Name] = &xmlstructs.Relationships{}
		}
	}

	return nil
}

func (e *state) loadXML(name string, target any) error {
	for _, f := range e.reader.File {
		if f.Name == name {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()
			return xml.NewDecoder(rc).Decode(target)
		}
	}
	return fmt.Errorf("file %s not found in zip", name)
}

func (e *state) writeXML(zw *zip.Writer, name string, data any) error {
	w, err := zw.Create(name)
	if err != nil {
		return err
	}
	fmt.Fprint(w, xml.Header)
	return xml.NewEncoder(w).Encode(data)
}

func (e *state) copyFile(f *zip.File, zw *zip.Writer) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	w, err := zw.Create(f.Name)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, rc)
	return err
}
