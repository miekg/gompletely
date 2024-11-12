package main

import (
	"bytes"
	"embed"
	"flag"
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
	zshtmpl  = Tmpl("zsh")

	flagShell = flag.String("s", "bash", "generate the completions for this shell.")
)

func main() {
	var (
		buf []byte
		err error
	)
	flag.Parse()

	if *flagShell != "bash" || *flagShell != "zsh" {
		log.Fatalf("invalid shell %q, expected %q or %q", *flagShell, "bash", "zsh")
	}

	if len(os.Args) == 1 {
		buf, err = io.ReadAll(os.Stdin)
	} else {
		if buf, err = os.ReadFile(os.Args[1]); err != nil {
			log.Fatalf("Can't read file %q: %s", os.Args[1], err)
		}
	}
	p := Patterns{}
	if err := yaml.Unmarshal(buf, &p); err != nil {
		log.Fatal(err)
	}

	if err := p.Valid(); err != nil {
		log.Fatal(err)
	}

	out := &bytes.Buffer{}
	filename := ""
	switch *flagShell {
	case "bash":
		b := p.Bash()
		if err = bashtmpl.Execute(out, b); err != nil {
			log.Fatal(err)
		}
		if len(os.Args) == 1 {
			fmt.Println(out.String())
			return
		}
		base := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1]))
		filename = base + ".bash"
	case "zsh":
		z := p.Zsh()
		if err = zshtmpl.Execute(out, z); err != nil {
			log.Fatal(err)
		}
		if len(os.Args) == 1 {
			fmt.Println(out.String())
			return
		}
		base := strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1]))
		filename = "_" + base
	}
	if err := os.WriteFile(filename, out.Bytes(), 0644); err != nil {
		log.Fatalf("Can't write file %q: %s", filename, err)
	}
}
