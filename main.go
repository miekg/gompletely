package main

import (
	"embed"
	"io"
	"os"

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
	if err = bashtmpl.Execute(os.Stdout, b); err != nil {
		log.Fatal(err)
	}
	return
}
