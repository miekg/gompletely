package main

import (
	"fmt"
	"slices"
	"strings"
)

type Zsh struct {
	Command  string
	Patterns []Arg
}

type Arg struct {
	Positional []string
	Options    []string // The completion it requires.
	// help is include in the above if there was one.
}

// Zsh returns a structure suitable for rendering in the zsh template.
func (p Patterns) Zsh() Zsh {
	z := Zsh{Command: p.Cmd()}
	keys := []string{}
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

	fmt.Printf("KEYS %+v\n", keys)

	args := []Arg{}
	for _, k := range keys {
		a := Arg{}
		actions := []string{}
		strs := []string{}
		println("KEEKE", k)
		for _, pat := range p[k] {
			switch pat.Type {
			case Command:
				// not which one
				a.Positional = append(a.Positional, pat.Completion)
				// need to come last
			case Option:
				a.Options = append(a.Options, pat.Completion)
			case Action:
				if pat.Completion != ActionNoop {
					actions = append(actions, "-A "+pat.Completion)
				}
			case String:
				strs = append(strs, pat.Completion)
			}

		}
		args = append(args, a)
	}
	z.Patterns = args
	return z
}
