package main

import (
	"fmt"
	"slices"
	"strings"
)

// In the zsh struct we rewrite the patterns we get from the yaml a bit, so that it can be used more easily in the
// template.
type Zsh struct {
	Command  string
	Commands []string // the sorted keys of the map Patterns.
	Patterns
}

// Zsh returns a structure suitable for rendering in the zsh template.
func (p Patterns) Zsh() Zsh {
	z := Zsh{Command: p.Cmd(), Patterns: map[string][]Pattern{}}
	z.Patterns = p

	keys := []string{}
	for k := range p {
		keys = append(keys, k)
	}
	// sort on key length, short -> less short.
	slices.SortFunc(keys, func(a, b string) int {
		ret := len(a) - len(b)
		if ret != 0 {
			return ret
		}
		return strings.Compare(a, b)
	})
	z.Commands = keys

	fmt.Printf("#compef _%s %s\n\n", z.Command, z.Command)
	for _, command := range z.Commands {
		patterns, ok := z.Patterns[command] // may be empty because we delete from the map
		if !ok {
			continue
		}
		fmt.Printf("function _%s {\n\tlocal line\n\n\t_arguments -C \\\n", command)

		// Options first
		for _, p := range patterns {
			if p.Type == Command {
				continue
			}
			if p.Help == "" {
				p.Help = "[]"
			}
			fmt.Printf("\t\t'%s%s", p.Completion, p.Help)
			args := z.Patterns.OptionHasArg(command, p.Completion)
			if args != nil {
				quoted := ""
				for i := range args {
					quoted += `"` + args[i] + `"`
				}

				fmt.Printf(":: _values %q %s" /*description*/, "userdb", quoted)
				// remove from pattersn
				delete(z.Patterns, command+"*"+p.Completion)
			}
			fmt.Printf("' \\\n")
		}
		fmt.Printf("\t\t\"*::arg:->args\"\n")
		fmt.Printf("}\n")
	}

	return z
}
