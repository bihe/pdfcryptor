package config

// PdfUtil defines the available PDF utility for encryption/decryption operations
type PdfUtil int

const (
	// QPDF utiliy http://qpdf.sourceforge.net/
	QPDF PdfUtil = iota
	// PDFTK utility https://www.pdflabs.com/tools/pdftk-the-pdf-toolkit/
	PDFTK PdfUtil = iota
)
