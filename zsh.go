package main

import (
	"bytes"
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

// Zsh returns a buffer with the completion. There is no template called.
func (p Patterns) Zsh() (Zsh, *bytes.Buffer) {
	b := &bytes.Buffer{}
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

	fmt.Fprintf(b, "#compdef %s\n", z.Command)
	fmt.Fprintf(b, "compdef _%s %s\n\n", z.Command, z.Command)
	fmt.Fprintf(b, "# Generated by gompletely (https://github.com/miekg/gompletely)\n\n")
	for _, command := range z.Commands {
		patterns, ok := z.Patterns[command] // may be empty because we delete from the map
		if !ok {
			continue
		}
		fmt.Fprintf(b, "function _%s {\n\tlocal line\n\n\t_arguments -C \\\n", funcName(command))

		// Options
		// --fs-type[]: : _values "userdb" zfs lvm dir' \
		for _, p := range patterns {
			if p.Position > 0 {
				continue
			}
			if p.Subcommand {
				continue
			}
			if p.Help == "" {
				p.Help = "[]"
			}

			args, help := z.Patterns.OptionHasArg(command, p.Completion)
			if args == nil {
				fmt.Fprintf(b, "\t\t'%s%s", p.Completion, p.Help)
			} else {
				fmt.Fprintf(b, "\t\t'%s%s", p.Completion, help)
				// The : : instead of :: is significant, between working _values, and not. It holds the description of what is being completed.
				if len(args) == 1 && p.Type == Action {
					fmt.Fprintf(b, ": : %s", strings.Join(args, " "))
				} else {
					fmt.Fprintf(b, ": : _values %q %s" /*description*/, "userdb", strings.Join(args, " "))
				}
				// remove from patterns, as we have handled it
				delete(z.Patterns, command+"*"+p.Completion)
			}
			fmt.Fprintf(b, "' \\\n")
		}

		// gather positional arguments with the same number, as they most be processs
		// on the same line in the _arguments ... Put those in a map[num]string, there should
		// only be a single one, per number: Check for this in Valid() TODO.
		poschoice := map[int]string{}
		for _, p := range patterns {
			if p.Position == 0 {
				continue
			}
			if p.Message == "" {
				continue
			}
			poschoice[p.Position] = p.Message
		}

		// check for subcommands
		subcommands := []string{}
		caseSub := ""
		for _, p := range patterns {
			if !p.Subcommand {
				continue
			}
			subcommands = append(subcommands, p.Message)
		}
		// create the stanza we need to add after the _arguments calling
		if len(subcommands) > 1 {
			caseSub = "\n\t\tcase $line[1] in\n"
			for _, s := range subcommands {
				cmd := "_" + strings.Replace(command, " ", "_", -1) + "_" + s
				caseSub += fmt.Sprintf("\t\t\t%s)\n\t\t\t\t%s\n\t\t\t;;\n", s, funcName(cmd))
			}
			caseSub += "\t\tesac\n"
		}

		// Positional arguments,
		//  "1: :(quietly loudly)" \
		for _, p := range patterns {
			if p.Subcommand && len(subcommands) > 0 {
				// subcommands default to position 1, handle them all
				fmt.Fprintf(b, "\t\t'1: :( %s )", strings.Join(subcommands, " "))
				fmt.Fprintf(b, "' \\\n")
				subcommands = []string{}
				continue
			}

			if p.Position == 0 {
				continue
			}
			if p.Subcommand {
				continue
			}

			if choice, ok := poschoice[p.Position]; ok {
				// if the completion is empty, this is mean as a hint to the user what to complete
				if p.Completion == "" {
					// 2:service name:'
					fmt.Fprintf(b, "\t\t'%d:%s:", p.Position, p.Message)
					fmt.Fprintf(b, "' \\\n")
				} else {
					fmt.Fprintf(b, "\t\t'%d: : _values %q ( %s )", p.Position /*description */, "userdb", choice)
					fmt.Fprintf(b, "' \\\n")
				}
				delete(poschoice, p.Position) // delete ourselves from the map
				continue
			}

			// if we are here wih a valid p.Message, we were deleted from the map above, skip
			if p.Message != "" {
				continue
			}
			if p.Type == Action {
				p.Completion = actionToZsh(p.Completion)
				fmt.Fprintf(b, "\t\t'%d: : %s", p.Position, p.Completion)
				fmt.Fprintf(b, "' \\\n")
				continue
			}

			fmt.Fprintf(b, "\t\t'%d: : _values %q %s", p.Position /*description */, "userdb", p.Completion)
			fmt.Fprintf(b, "' \\\n")
		}

		fmt.Fprintf(b, "\t\t\"*::arg:->args\"\n")

		if caseSub != "" {
			fmt.Fprintf(b, caseSub)
		}

		fmt.Fprintf(b, "}\n")
	}
	return z, b
}

// funcName returns a string the valid function name in Zsh.
func funcName(cmd string) string {
	s1 := strings.Replace(cmd, " ", "_", -1)
	s1 = strings.Replace(s1, "-", "_", -1)
	return s1
}

// actionToZsh
func actionToZsh(a string) string {
	switch a {
	case "file", "directory":
		return "_files"
	case "group":
		return "_groups"
	case "user":
		return "_users"
	case "export":
		return "_parameters"
	}

	return ""
}
