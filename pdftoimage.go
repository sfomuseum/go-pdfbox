package pdfbox

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// type PDFToImageCallback is a callback function invoked by the PDFToImage method. It passes the path and an
// `io.Reader` instance for each image produced by the PDFToImage method.
type PDFToImageCallback func(context.Context, string, io.Reader) error

// PDFToImage() is an utility method for invoking the pdfbox `PDFToImage` tool for a file derived from 'r'. This method only
// exposes the "-startPage" and "-endPage" parameters. Resultant images are passed to 'cb' as `io.Reader` instances.
func PDFToImage(ctx context.Context, pdfb *PDFBox, r io.Reader, start_page int, end_page int, cb PDFToImageCallback) error {

	tmpfile_r, err := ioutil.TempFile("", "pdfbox-images.*.pdf")

	if err != nil {
		return fmt.Errorf("Failed to create tempfile for reader, %w", err)
	}

	defer os.Remove(tmpfile_r.Name())

	_, err = io.Copy(tmpfile_r, r)

	if err != nil {
		return fmt.Errorf("Failed to write tempfile for reader, %w", err)
	}

	err = tmpfile_r.Close()

	if err != nil {
		return fmt.Errorf("Failed to close tempfile for reader, %w", err)
	}

	tmpdir, err := ioutil.TempDir("", "pdfbox-images")

	if err != nil {
		return fmt.Errorf("Failed to create tempdir for writing PDF images, %w", err)
	}

	defer os.RemoveAll(tmpdir)

	fname := filepath.Base(tmpfile_r.Name())
	ext := filepath.Ext(fname)

	prefix := strings.Replace(fname, ext, "-", 1)
	abs_prefix := filepath.Join(tmpdir, prefix)

	str_start := strconv.Itoa(start_page)
	str_end := strconv.Itoa(end_page)

	err = pdfb.Execute(ctx, "PDFToImage", "-startPage", str_start, "-endPage", str_end, "-outputPrefix", abs_prefix, tmpfile_r.Name())

	if err != nil {
		return fmt.Errorf("Failed to extract images, %w", err)
	}

	tmpfs := os.DirFS(tmpdir)

	err = fs.WalkDir(tmpfs, ".", func(path string, d fs.DirEntry, err error) error {

		if !strings.HasPrefix(path, prefix) {
			return nil
		}

		r, err := tmpfs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open '%s', %w", path, err)
		}

		defer r.Close()

		err = cb(ctx, path, r)

		if err != nil {
			return fmt.Errorf("Image callback for %s failed, %w", path, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to walk tempdir, %w", err)
	}

	return nil
}
