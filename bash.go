package main

import (
	"bytes"
	"fmt"
	"slices"
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
	Completion string // The compgen to add.
	Positional string // positional argument switch (only used in "*"-case" and if there are positional arguments)
}

// Bash returns a structure suitable for rendering in the bash template.
func (p Patterns) Bash() Bash {
	b := Bash{Command: p.Cmd()}
	keys := []string{}
	for k := range p {
		keys = append(keys, k)
	}
	// sort on key length, shortest ones need to be at the end for the case to work correctly.
	slices.SortFunc(keys, func(a, b string) int {
		ret := len(b) - len(a)
		if ret != 0 {
			return ret
		}
		return strings.Compare(a, b)
	})

	pos := []Case{}
	// The b.Command key pattern is for the toplevel command. For this command we _also_ inject positional
	// argument completion. We grab the toplevel, Action, Command and Strings. If _more_ than one inject this
	i := 1
	for _, pat := range p[b.Command] {
		if pat.Type == Option || pat.Type == String { // only do command and actions
			continue
		}
		c := Case{CaseString: quote(strconv.FormatInt(int64(pat.Position), 10))}
		switch pat.Type {
		case Command:
			c.Completion = fmt.Sprintf(`-W "$(_%s_completions_filter "%s")"`, b.Command, pat.Completion)
		case Action:
			c.Completion = "-A " + pat.Completion
		}
		pos = append(pos, c)
		i++
	}

	// Only when we have 2 or more positional arguments will we need to fill the extra switch. Fill
	// out the template, for later use, when builing out the cases below.
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
	c := Case{}
	for _, k := range keys {
		casestring := strings.TrimPrefix(k, b.Command)
		fields := strings.Split(casestring, "*")
		switch len(fields) {
		case 0:
			c.CaseString = "*"
		case 1:
			c.CaseString = quote(fields[0]) + "*"
		case 2:
			if fields[0] == "" {
				c.CaseString = "*" + quote(fields[1])
			} else {
				c.CaseString = quote(fields[0]) + "*" + quote(fields[1])
			}
		default:
			for i := range fields {
				fields[i] = quote(fields[i])
			}
			c.CaseString = strings.Join(fields, "*")
		}

		commands := []string{}
		options := []string{}
		actions := []string{}
		strs := []string{}
		subcommands := []string{}
		for _, pat := range p[k] {
			if pat.Subcommand {
				subcommands = append(subcommands, pat.Message)
				continue
			}

			switch pat.Type {
			case Command:
				commands = append(commands, pat.Completion)
			case Option:
				options = append(options, pat.Completion)
			case Action:
				actions = append(actions, "-A "+pat.Completion)
			case String:
				strs = append(strs, pat.Completion)
			}
		}

		completions_filter := strings.TrimSpace(join(options) + join(commands) + join(strs) + join(subcommands))
		switch completions_filter {
		case "":
			if len(actions) > 0 {
				c.Completion = fmt.Sprintf(`%s`, strings.TrimSpace(join(actions)))
			}
		default:
			c.Completion = fmt.Sprintf(`%s-W "$(_%s_completions_filter "%s")"`, join(actions), b.Command, completions_filter)
		}
		if c.CaseString == "*" { // inject massive switch
			c.Positional = posbuf.String()
		}
		if c.Completion != "" {
			patterns = append(patterns, c)
		}

	}
	b.Patterns = patterns
	return b
}

const postmpl = `
	COMP_CARG=$COMP_CWORD; for i in "${COMP_WORDS[@]}"; do [[ ${i} == -* ]] && ((COMP_CARG = COMP_CARG - 1)); done
	case $COMP_CARG in
	{{range .}}
	  {{.CaseString}})
	    while read -r; do COMPREPLY+=("$REPLY"); done < <(compgen {{.Completion}} -- "$cur")
            return
          ;;{{end}}
        esac
`
