# Thoth Document Processing Library

A modern, extensible Go library (Go 1.26+) for high-performance processing of Microsoft Excel (.xlsx), Microsoft Word (.docx), and Adobe Acrobat PDF (.pdf) documents.

## Key Features

- **Multi-Format Support**: Single library to handle the most common document formats.
- **Fluent API**: Chainable, intuitive API for document manipulation.
- **Advanced Styling**: Granular control over fonts, colors, borders, and alignment using a builder pattern.
- **Table Management**: Systematic API for creating and managing complex tables with cell merging.
- **Infrastructure Ready**: Built-in support for S3 storage and URL-based document loading.
- **Performance Focused**: Efficient memory usage with `sync.Pool` and streaming capabilities.
- **Structured Logging**: Integrated with `zerolog` for enterprise-grade observability.
- **Modern Go**: Leverages Go 1.26+ features.

## Format-Specific Features

### 📊 Excel (.xlsx)
- Create new spreadsheets or open existing ones.
- Scoped sheet handles for clean, chainable operations.
- Cell value setting (strings, numbers, dates).
- Advanced styling: bold, italic, colors, borders, and number formats.
- Column width management and cell merging.
- AutoFilter and Freeze Panes.
- **Image insertion** into worksheets.

### 📝 Word (.docx)
- Paragraph-based content addition with full styling support.
- **Image insertion** with positioning.
- Page breaks and section management.
- Complex table API with row/cell scoping and cell merging.
- Document metadata management.

### 📄 PDF (.pdf)
- High-fidelity PDF generation with **Stream Compression** (FlateDecode).
- Page settings: orientation (portrait/landscape) and paper types (A4, Letter, etc.).
- Headers and footers with automatic repetition and page numbering (`{n}` of `{nb}`).
- **Navigation**: Automatic **Bookmarks (Outlines)** from Title/Heading styles.
- **Interactivity**: **Hyperlinks** (External URLs) support in text paragraphs.
- **Shape drawing** (Lines, Rectangles).
- **Image insertion** into document flow.
- Table support similar to the Word API.

### ☁️ Storage & Core
- **S3 Integration**: Open documents directly from Amazon S3.
- **URL Support**: Fetch documents from HTTP/HTTPS endpoints.
- **Context Awareness**: All I/O operations respect `context.Context` for timeouts and cancellation.
- **Plugable Storage**: Interface-based storage providers.

## Installation

```bash
go get github.com/gsoultan/thoth
```

## Quick Start

```go
import (
    "context"
    "github.com/gsoultan/thoth/core"
    "github.com/rs/zerolog"
    "os"
)

func main() {
    logger := zerolog.New(os.Stdout)
    thoth := core.New(logger)
    ctx := context.Background()

    // Create a new Excel document
    ss, _ := thoth.Excel().New(ctx)
    defer ss.Close()

    sheet, _ := ss.Sheet("Main")
    sheet.Cell("A1").Set("Hello from Thoth!")
    
    ss.Export("example.xlsx")
}
```

## Examples

The `examples/` directory contains comprehensive usage samples:
- `examples/main.go`: General overview of Excel, Word, and PDF features.
- `examples/pdf_complex/main.go`: Advanced PDF features including headers, footers, and shapes.
- `examples/pdf_production/main.go`: Production-ready PDF features: compression, bookmarks, hyperlinks, and pagination.

## Project Structure

- `core/`: High-level entry points and infrastructure (Storage, Logging).
- `excel/`: Domain-specific logic for Excel.
- `word/`: Domain-specific logic for Word.
- `pdf/`: Domain-specific logic for PDF.
- `document/`: Shared interfaces and common types.

## License
MIT