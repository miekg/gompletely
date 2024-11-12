package main

import (
	"slices"
	"strings"
)

type Zsh struct {
	Command string
	Args    []Arg
}

type Arg struct {
	Option     string // The case string to switch on.
	Completion string // The completion it requires.
	Help       string // The help string between [ and ].
	Positional string // positional argument.
}

// Zsh returns a structure suitable for rendering in the zsh template.
func (p Patterns) Zsh() Bash {
	b := Zsh{Command: p.Cmd()}
	for k := range p {
		keys = append(keys, k)
	}
	// sort on key length, sortest ones need to be at the end for the case to work correctly.
	slices.SortFunc(keys, func(a, b string) int {
		ret := len(b) - len(a)
		if ret != 0 {
			return ret
		}
		return strings.Compare(a, b)
	})

	for _, pat := range p[b.Command] {
	}
}
