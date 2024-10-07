package main

import (
	"fmt"
	"sort"
	"strings"
)

type Bash struct {
	Command  string
	Patterns []Case
}

type Case struct {
	CaseString string // The case string to switch on.
	CompGen    string // The compgen to add.
}

// ToBash returns a structure suitable for rendering in the template.
func ToBash(p Patterns) Bash {
	b := Bash{Command: Cmd(p)}
	keys := []string{}
	for k := range p {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

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
		for _, p := range p[k] {
			switch p.CompType {
			case Command:
				commands = append(commands, p.CompGen)
			case Option:
				options = append(options, p.CompGen)
			case Action:
				actions = append(actions, "-A "+p.CompGen)
			case String:
				strs = append(strs, p.CompGen)
			}
		}

		compgen := fmt.Sprintf(`%s-W "$(_%s_completions_filter "%s")"`,
			join(actions),
			b.Command, strings.TrimSpace(
				join(options)+join(commands)+join(strs),
			),
		)
		c.CompGen = compgen

		patterns = append(patterns, c)
	}
	b.Patterns = patterns
	return b
}
