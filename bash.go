package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type Bash struct {
	Command  string
	Patterns []Case
}

type Case struct {
	CaseString string // The case string to switch on.
	CompGen    string // The compgen to add.
	Positional string // positional switch (only used for "*"-case")
}

// ToBash returns a structure suitable for rendering in the template.
func ToBash(p Patterns) Bash {
	b := Bash{Command: Cmd(p)}
	keys := []string{}
	for k := range p {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	pos := []Case{}
	// The empty key pattern is for the toplevel command. For this command we _also_ inject positional
	// argument completion. Grab the toplevel, Action, Command and String. If _more_ than one inject this
	i := 1
	for _, pat := range p[b.Command] {
		if pat.CompType == Option {
			continue
		}
		c := Case{CaseString: quote(strconv.FormatInt(int64(i), 10))}
		switch pat.CompType {
		case Command:
			c.CompGen = fmt.Sprintf(`-W "$(_%s_completions_filter "%s")"`, b.Command, pat.CompGen)
		case Action:
			c.CompGen = "-A " + pat.CompGen
		case String:
			c.CompGen = fmt.Sprintf(`-W "$(_%s_completions_filter "%s")"`, b.Command, pat.CompGen)
		}

		pos = append(pos, c)
		i++
	}

	// Only when we have 2 or more positional arguments will we need to fill the extra switch. Fill
	// out the template, for later use.
	posbuf := &bytes.Buffer{}
	if len(pos) > 1 {
		tmpl, err := template.New("test").Parse(postmpl)
		if err != nil {
			panic("Invalid postmpl: " + err.Error())
		}
		if err := tmpl.Execute(posbuf, pos); err != nil {
			panic("Invalid postmpl: " + err.Error())
		}
	}

	patterns := []Case{}
	for _, k := range keys {
		casestring := strings.TrimPrefix(k, b.Command)
		switch {
		case strings.HasPrefix(casestring, " "):
			casestring = quote(strings.TrimSpace(casestring)) + "*"
		case strings.HasPrefix(casestring, "*"):
			casestring = "*" + quote(casestring[1:])
		case casestring == "":
			casestring = "*"
		}

		c := Case{CaseString: casestring}
		commands := []string{}
		options := []string{}
		actions := []string{}
		strs := []string{}
		for _, pat := range p[k] {
			switch pat.CompType {
			case Command:
				commands = append(commands, pat.CompGen)
			case Option:
				options = append(options, pat.CompGen)
			case Action:
				actions = append(actions, "-A "+pat.CompGen)
			case String:
				strs = append(strs, pat.CompGen)
			}
		}

		compgen := fmt.Sprintf(`%s-W "$(_%s_completions_filter "%s")"`,
			join(actions),
			b.Command, strings.TrimSpace(
				join(options)+join(commands)+join(strs),
			),
		)
		c.CompGen = compgen
		if c.CaseString == "*" {
			c.Positional = posbuf.String()
		}
		patterns = append(patterns, c)

	}
	b.Patterns = patterns
	return b
}

const postmpl = `
	COMP_CARG=$COMP_CWORD; for i in "${COMP_WORDS[@]}"; do [[ ${i} == -* ]] && ((COMP_CARG = COMP_CARG - 1)); done
	case $COMP_CARG in
	{{range .}}
	  {{.CaseString}})
	    while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen {{.CompGen}} -- "$cur")
            return
          ;;{{end}}
        esac
`
