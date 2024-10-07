package main

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.science.ru.nl/log"
	"gopkg.in/yaml.v3"
)

//go:embed *.go.tmpl
var tmplfs embed.FS

var (
	bashtmpl = Tmpl("bash")
	// zshtmpl = Tmpl("zsh")
)

func main() {
	var (
		buf []byte
		err error
	)

	if len(os.Args) == 1 {
		buf, err = io.ReadAll(os.Stdin)
	} else {
		if buf, err = os.ReadFile(os.Args[1]); err != nil {
			log.Fatalf("Can't read file %q: %s", os.Args[1], err)
		}
	}
	p := Patterns{}
	if err := yaml.Unmarshal(buf, &p); err != nil {
		log.Error(err)
	}

	b := ToBash(p)
	out := &bytes.Buffer{}
	if err = bashtmpl.Execute(out, b); err != nil {
		log.Fatal(err)
	}
	if len(os.Args) == 1 {
		fmt.Println(out.String())
		return
	}
	base := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1]))
	if err := os.WriteFile(base+".bash", out.Bytes(), 0644); err != nil {
		log.Fatalf("Can't write file %q: %s", base, err)
	}
}
