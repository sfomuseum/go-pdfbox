package pdfbox

import (
	"context"
	"os/exec"
	"fmt"
	"io"
	"io/ioutil"
)

type PDFBox struct {
	java string
	jar string
}

func New() (*PDFBox, error){

	p := &PDFBox{}
	return p, nil
}

func (p *PDFBox) Execute(ctx context.Context, command string, args ...interface{}) error {

	if len(args) == 0 {
		return fmt.Errorf("No arguments")
	}
		
	cmd := exec.CommandContext(p.java, "-jar", p.jar, command, args...)
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("Failed to execute %s, %w", command, err)		
	}

	return nil
}
