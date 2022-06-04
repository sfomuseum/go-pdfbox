package pdfbox

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-pdfbox/jar"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type PDFBox struct {
	java    string
	jarfile string
}

func New(ctx context.Context, uri string) (*PDFBox, error) {

	// check for java here, use sync.Once

	jar_r, err := jar.FS.Open("pdfbox-app-2.0.26.jar")

	if err != nil {
		return nil, fmt.Errorf("Failed to open jar file, %w", err)
	}

	jar_wr, err := ioutil.TempFile("", "pdfbox")

	if err != nil {
		return nil, fmt.Errorf("Failed to create tempfile for jar, %w", err)
	}

	_, err = io.Copy(jar_wr, jar_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to copy jarfile, %w", err)
	}

	err = jar_wr.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to close (temp) jarfile, %w", err)
	}

	p := &PDFBox{
		jarfile: jar_wr.Name(),
	}

	return p, nil
}

func (p *PDFBox) Close() error {
	return os.Remove(p.jarfile)
}

func (p *PDFBox) ExecuteWithReader(ctx context.Context, r io.Reader, command string, args ...interface{}) error {

	return nil
}

func (p *PDFBox) ExecuteWithReaderAndWriter(ctx context.Context, r io.Reader, wr io.Writer, command string, args ...interface{}) error {

	return nil
}

func (p *PDFBox) Execute(ctx context.Context, command string, args ...string) error {

	if len(args) == 0 {
		return fmt.Errorf("No arguments")
	}

	local_args := []string{
		"-jar",
		p.jarfile,
		command,
	}

	for _, a := range args {
		local_args = append(local_args, a)
	}

	cmd := exec.CommandContext(ctx, p.java, local_args...)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to execute %s, %w", command, err)
	}

	return nil
}
