package main

import (
	"fmt"
	"sort"
	"strconv"
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

	postional := []Case{}
	// The empty key pattern is for the toplevel command. For this command we _also_ inject positional
	// argument completion.
	// grab the toplevel, Action, Command and String. If _more_ than one inject this
	i := 2
	for _, pat := range p[b.Command] {
		if pat.CompType == Option {
			continue
		}
		println(pat.CompGen)
		println(pat.CompType)
		if pat.CompType == Action && pat.CompGen == ActionNone {
			i++
			continue
		}
		c := Case{CaseString: strconv.FormatInt(int64(i), 10)}
		switch pat.CompType {
		case Command:
			c.CompGen = fmt.Sprintf(`-W "$(_%s_completions_filter "%s")"`, b.Command, pat.CompGen)
		case Action:
			c.CompGen = "-A " + pat.CompGen
		case String:
			c.CompGen = fmt.Sprintf(`-W "$(_%s_completions_filter "%s")"`, b.Command, pat.CompGen)
		}

		postional = append(postional, c)
		i++
	}

	fmt.Printf("%+v\n", postional)

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
				if pat.CompGen == ActionNone {
					continue
				}
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

		patterns = append(patterns, c)
	}
	b.Patterns = patterns
	return b
}

const carg = `COMP_CARG=$COMP_CWORD; for i in "${COMP_WORDS[@]}"; do [[ ${i} == -* ]] && ((COMP_CARG = COMP_CARG - 1)); done`
