package pdfbox

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sfomuseum/go-pdfbox/jar"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

const READER string = "{READER}"

const WRITER string = "{WRITER}"

// type PDFBox is a struct for executing `pdfbox` commands.
type PDFBox struct {
	// The path to the locally installed Java binary.
	java string
	// The path to the temporary pdfbox jar file used by `PDFBox`
	jarfile string
}

// whichJava() attempts to determine that path of the locally installed Java binary using
// the `which` utility. This will probably not work on Windows and some versions of Unix.
func whichJava(ctx context.Context) (string, error) {

	// This will not work on Windows and some versions of Unix

	cmd := exec.CommandContext(ctx, "which", "java")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("Failed to determine which java, %w", err)
	}

	output = bytes.TrimSpace(output)
	path_java := string(output)

	_, err = os.Stat(path_java)

	if err != nil {
		return "", fmt.Errorf("Failed to stat '%s', %w", path_java, err)
	}

	return path_java, nil
}

// New() returns a new `PDFBox` instance configured by 'uri' which is expected to take
// the form of:
//
//	pdfbox://
func New(ctx context.Context, uri string) (*PDFBox, error) {

	// something something something sync.Once

	java, err := whichJava(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to locate java, %w", err)
	}

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
		java:    java,
		jarfile: jar_wr.Name(),
	}

	return p, nil
}

// Close() ...
func (p *PDFBox) Close() error {
	return os.Remove(p.jarfile)
}

// ExecuteWithReader() invokes pdfbox using 'command' and 'args' and the value of 'r' (written to a temporary file) as its input.
func (p *PDFBox) ExecuteWithReader(ctx context.Context, r io.Reader, command string, args ...string) error {

	tmpfile_r, err := ioutil.TempFile("", "pdfbox")

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

	set_reader := false

	for idx, a := range args {

		if a == READER {
			args[idx] = tmpfile_r.Name()
			set_reader = true
			break
		}
	}

	if !set_reader {
		return fmt.Errorf("Failed to set reader path")
	}

	err = p.Execute(ctx, command, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute %s, %w", command, err)
	}

	return nil
}

// ExecuteWithReaderAndWrite() invokes pdfbox using 'command' and 'args' and the value of 'r' (written to a temporary file) as its input, writing the final output to 'wr'.
func (p *PDFBox) ExecuteWithReaderAndWriter(ctx context.Context, r io.Reader, wr io.Writer, command string, args ...string) error {

	tmpfile_r, err := ioutil.TempFile("", "pdfbox")

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

	tmpfile_wr, err := ioutil.TempFile("", "pdfbox")

	if err != nil {
		return fmt.Errorf("Failed to create tempfile for writer, %w", err)
	}

	defer os.Remove(tmpfile_wr.Name())

	err = tmpfile_r.Close()

	if err != nil {
		return fmt.Errorf("Failed to close tempfile for writer, %w", err)
	}

	set_reader := false
	set_writer := false

	for idx, a := range args {

		if a == READER {
			args[idx] = tmpfile_r.Name()
			set_reader = true
			break
		}

		if a == WRITER {
			args[idx] = tmpfile_wr.Name()
			set_writer = true
			break
		}

	}

	if !set_reader {
		return fmt.Errorf("Failed to set reader path")
	}

	if !set_writer {
		return fmt.Errorf("Failed to set writer path")
	}

	err = p.Execute(ctx, command, args...)

	if err != nil {
		return fmt.Errorf("Failed to execute %s, %w", command, err)
	}

	wr_r, err := os.Open(tmpfile_wr.Name())

	if err != nil {
		return fmt.Errorf("Failed to open tempfile for writer, %w", err)
	}

	_, err = io.Copy(wr, wr_r)

	if err != nil {
		return fmt.Errorf("Failed to copy tempfile for writer to writer, %w", err)
	}

	return nil
}

// Execute() invokes pdfbox using 'command' and 'args' as the command input.
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

	// fmt.Println("DEBUG", p.java, local_args)

	cmd := exec.CommandContext(ctx, p.java, local_args...)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Failed to execute %s, %w (%s)", command, err, string(out))
	}

	return nil
}
