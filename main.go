package main

import (
	"embed"
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
	buf, err := os.ReadFile("testdata/AddVolume.yml")
	if err != nil {
		log.Fatal(err)
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
