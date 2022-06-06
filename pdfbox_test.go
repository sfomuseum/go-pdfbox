package pdfbox

import (
	"context"
	"testing"
)

func TestPDFBox(t *testing.T) {

	ctx := context.Background()

	p, err := New(ctx, "pdfbox://")

	if err != nil {
		t.Fatalf("Failed to create PDFBox, %v", err)
	}

	// create PDF from text

	// extract text from PDF

	err = p.Close()

	if err != nil {
		t.Fatalf("Failed to close PDFBox, %v", err)
	}
}

func TestPDFBoxExecute(t *testing.T) {

	t.Skip()
}

func TestPDFBoxExecuteWithReader(t *testing.T) {

	t.Skip()
}

func TestPDFBoxExecuteWithWriter(t *testing.T) {

	t.Skip()
}
