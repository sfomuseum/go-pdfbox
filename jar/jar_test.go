package jar

import (
	"io"
	"testing"
)

func TestFS(t *testing.T) {

	jarfile := "pdfbox-app-2.0.26.jar"
	fsize := int64(10081077)

	r, err := FS.Open(jarfile)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", jarfile, err)
	}

	i, err := io.Copy(io.Discard, r)

	if err != nil {
		t.Fatalf("Failed to copy %s, %v", jarfile, err)
	}

	if i != fsize {
		t.Fatalf("Invalid filesize for %s: %d", jarfile, i)
	}

}
