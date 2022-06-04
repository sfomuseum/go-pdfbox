package jar

import (
	"embed"
)

//go:embed *.jar
var FS embed.FS
