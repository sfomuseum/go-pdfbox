# go-pdfbox

Go package for working with the `pdfbox` command line tool, with an embedded copy of the pdfbox jarfile.

## Motivation

This is a simple Go "wrapper" package for working with the `pdfbox` command line tool. It contains an embedded copy of the pdfbox jar file but depends on a local copy of Java to run. The goal of the tool is to hide some of the details of working with pdfbox in a Go context but my hope is that eventually this package can be retired in favour of being able to invoke a WASM-compiled version of `pdfbox`. Until then this package exists.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-pdfbox.svg)](https://pkg.go.dev/github.com/sfomuseum/go-pdfbox)

## Example

_Error handling has been removed for the sake of brevity._

### Basic

```
import (
	"context"       
	"github.com/sfomuseum/go-pdfbox"
)

func main() {

     ctx := context.Background()
     p, _ := pdfbox.New(ctx, "pdfbox://")
     o.Execute(ctx, "ExtractText", "example.pdf", "example.txt")
}     
```

### Use with io.Reader and io.Writer

```
import (
       	"bytes"
	"bufio"		
	"context"       
	"github.com/sfomuseum/go-pdfbox"
	"os"
)

func main() {

     ctx := context.Background()

     r, _ := os.Open("example.pdf")

     var buf bytes.Buffer
     wr := bufio.NewWriter(&buf)

     p, _ := pdfbox.New(ctx, "pdfbox://")
     
     p.ExecuteWithReaderAndWriter(ctx, r, wr, "ExtractText", pdfbox.READER, pdfbox.WRITER)
}     
```

Note the use of the `pdfbox.READER` and `pdfbox.WRITER` variables. They are placeholder strings used to swap in the name of the temporary files created using the `io.Reader` and `io.Writer` variables respectively. The order of input and output files in pdfbox is not constant so it's easier just to be explicit about things rather than try to guess at positional elements.

## See also

* https://pdfbox.apache.org/
