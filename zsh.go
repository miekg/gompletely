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

	// possible fixing of patterns
	for _, command := range z.Commands {
		command = command
	}

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

	fmt.Printf("%s\n", z.Command)
	for _, command := range z.Commands {
		fmt.Printf("%s\n%+v\n", command, z.Patterns[command])
	}

	return z
}
