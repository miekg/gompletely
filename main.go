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

	flagShell   = flag.String("s", "bash", "generate the completions for this shell.")
	flagMessage = flag.String("m", "", "optional message to use when completing subcommands.")
)

func main() {
	var (
		buf []byte
		err error
	)
	flag.Parse()

	if *flagShell != "bash" && *flagShell != "zsh" {
		log.Fatalf("invalid shell %q, expected %q or %q", *flagShell, "bash", "zsh")
	}

	if flag.NArg() == 0 {
		buf, err = io.ReadAll(os.Stdin)
	} else {
		if buf, err = os.ReadFile(flag.Arg(0)); err != nil {
			log.Fatalf("Can't read file %q: %s", flag.Arg(0), err)
		}
	}
	p := Patterns{}
	if err := yaml.Unmarshal(buf, &p); err != nil {
		log.Fatalf("can't parse %s: %s", flag.Arg(0), err)
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

		if flag.NArg() == 0 {
			fmt.Println(out.String())
			return
		}
		base := strings.TrimSuffix(flag.Arg(0), filepath.Ext(flag.Arg(0)))
		filename = base + ".bash"
	case "zsh":
		_, out = p.Zsh()
		if flag.NArg() == 0 {
			fmt.Print(out.String())
			return
		}
		base := strings.TrimSuffix(flag.Arg(0), filepath.Ext(flag.Arg(0)))
		filename = "_" + base
	}
	if err := os.WriteFile(filename, out.Bytes(), 0644); err != nil {
		log.Fatalf("Can't write file %q: %s", filename, err)
	}
}
