package pdfbox

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
	"testing"
)

func TestPDFToImage(t *testing.T) {

	ctx := context.Background()

	test_pdf := "fixtures/helloworld.pdf"

	r, err := os.Open(test_pdf)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", test_pdf, err)
	}

	defer r.Close()

	pdfb, err := New(ctx, "pdfbox://")

	if err != nil {
		t.Fatalf("Failed to create PDFBox, %v", err)
	}

	cb := func(ctx context.Context, path string, r io.Reader) error {

		_, _, err := image.Decode(r)

		if err != nil {
			return fmt.Errorf("Failed to decode %s, %w", path, err)
		}

		return nil
	}

	err = PDFToImage(ctx, pdfb, r, 1, 1, cb)

	if err != nil {
		t.Fatalf("Failed to convert to PDF to image, %v", err)
	}

}
