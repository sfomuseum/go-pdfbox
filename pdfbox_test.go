package pdfbox

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
)

func TestPDFBox(t *testing.T) {

	ctx := context.Background()

	p, err := New(ctx, "pdfbox://")

	if err != nil {
		t.Fatalf("Failed to create PDFBox, %v", err)
	}

	// Create PDF from text

	tmp_txt, err := ioutil.TempFile("", "pdfbox-testing")

	if err != nil {
		t.Fatalf("Failed to create text tempfile, %v", err)
	}

	defer os.Remove(tmp_txt.Name())

	_, err = tmp_txt.Write([]byte("Hello world"))

	if err != nil {
		t.Fatalf("Failed to write text tempfile, %v", err)
	}

	err = tmp_txt.Close()

	if err != nil {
		t.Fatalf("Failed to close text tempfile, %v", err)
	}

	tmp_pdf, err := ioutil.TempFile("", "pdfbox-testing")

	if err != nil {
		t.Fatalf("Failed to create pdf tempfile, %v", err)
	}

	defer os.Remove(tmp_pdf.Name())

	err = tmp_pdf.Close()

	if err != nil {
		t.Fatalf("Failed to close pdf tempfile, %v", err)
	}

	err = p.Execute(ctx, "TextToPDF", tmp_pdf.Name(), tmp_txt.Name())

	if err != nil {
		t.Fatalf("Failed to convert text to PDF, %v", err)
	}

	// Extract text from PDF

	// Clean up PDFBox

	err = p.Close()

	if err != nil {
		t.Fatalf("Failed to close PDFBox, %v", err)
	}

}
