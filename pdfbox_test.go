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
