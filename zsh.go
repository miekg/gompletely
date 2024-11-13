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

	/* werkend
	'--fs-type[fs type2]: : _values "userdb" zfs lvm dir' \
	*/

	fmt.Printf("#compdef _%s %s\n\n", z.Command, z.Command)
	for _, command := range z.Commands {
		patterns, ok := z.Patterns[command] // may be empty because we delete from the map
		if !ok {
			continue
		}
		fmt.Printf("function _%s {\n\tlocal line\n\n\t_arguments -C \\\n", funcName(command))

		// Options
		// --fs-type[]: : _values "userdb" zfs lvm dir' \
		for _, p := range patterns {
			if p.Position > 0 {
				continue
			}
			if p.Help == "" {
				p.Help = "[]"
			}
			fmt.Printf("\t\t'%s%s", p.Completion, p.Help)
			args := z.Patterns.OptionHasArg(command, p.Completion)
			if args != nil {
				// the : : instead of :: is significant, between working _values, and not.
				// It holds the description of what is being completed.
				fmt.Printf(": : _values %q %s" /*description*/, "userdb", strings.Join(args, " "))
				// remove from patterns, as we have handled it
				delete(z.Patterns, command+"*"+p.Completion)
			}
			fmt.Printf("' \\\n")
		}

		// gather positional arguments with the same number, as they most be processs
		// on the same line in the _arguments ... Put those in a map[num][]string
		poschoice := map[int][]string{}
		for _, p := range patterns {
			if p.Position == 0 {
				continue
			}
			if p.PosChoice == "" {
				continue
			}
			poschoice[p.Position] = append(poschoice[p.Position], p.PosChoice)
			poschoice[p.Position] = slices.Compact(poschoice[p.Position])
		}

		// Positional arguments,
		//  "1: :(quietly loudly)" \
		for _, p := range patterns {
			if p.Position == 0 {
				continue
			}
			if choices, ok := poschoice[p.Position]; ok {
				fmt.Printf("\t\t'%d: : _values %q ( %s )", p.Position /*description */, "userdb", strings.Join(choices, " "))
				fmt.Printf("' \\\n")
				delete(poschoice, p.Position) // delete ourselves from the map
				continue
			}
			// if we are here wih a valid p.PosChoice, we were deleted from the map above, skip
			if p.PosChoice != "" {
				continue
			}
			if p.Type == Command {
				comp := strings.TrimPrefix(p.Completion, "$(")
				comp = strings.TrimSuffix(comp, ")")
				p.Completion = "{ " + comp + " }"
			}

			fmt.Printf("\t\t'%d: : _values %q %s", p.Position /*description */, "userdb", p.Completion)
			fmt.Printf("' \\\n")
		}

		fmt.Printf("\t\t\"*::arg:->args\"\n")

		// TODO: subcommands, and correctly generate those functions.
		fmt.Printf("}\n")
	}
	return z
}

// funcName returns a string the valid function name in Zsh.
func funcName(cmd string) string {
	return strings.Replace(cmd, " ", "_", -1)
}
