package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Patterns is a whole file of completions for a command.
type Patterns map[string][]Pattern

// OptionHasArg search the keys in p with <cmd>*option, and returns the completions
func (p Patterns) OptionHasArg(cmd, option string) []string {
	patterns, ok := p[cmd+"*"+option]
	if !ok {
		return nil
	}
	// ok this is an option, return a string valid for _values, concattenate them with spaces
	// and quota them
	cs := make([]string, len(patterns))
	for i := range patterns {
		if patterns[i].Type == Action {
			cs[i] = actionToZsh(patterns[i].Completion)
			continue
		}
		cs[i] = patterns[i].Completion
	}
	return cs
}

// A type has four (useful) values:
// An Action is a bash shell action to e.g. complete files. These are denoted in the yaml as "<action>".
// There is one special action "<none>" which is used for arguments that don't have any completion.
//
// An Option is the option that can be completed for a command, i.e. "-f" or "--force". These are recognized because
// they all start with a minus,
//
// Command is a shell command to be executed, these are recognized because they use the command substituion
// syntax "$(echo hello)"
//
// String is a single string that is the completion string itself.
//
// Optionally these values may be followed by a "[help]" string in blockquotes that can be used as a help. In the future
// this might be following again with ":message" which zsh uses to say what is being completed.
type CompletionType int

const (
	None CompletionType = iota
	Action
	Option
	Command
	String
)

const ActionNoop = "noop"

// Pattern is a completion we read from the yaml. It is altered and made suitable for completion generation by Bash/Zsh/... etc.
type Pattern struct {
	Type       CompletionType
	Completion string // the literal completion string
	Position   int    // if > 0 this is a positional argument
	PosChoice  string // if Poistion > 0 , but there are several options, we use PosChose string to differentiate
	Help       string // optional help text
}

func (p *Pattern) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.Decode(&str)
	if err != nil {
		return err
	}

	// --root[bla] --> [bla] --root
	help, str := stripHelp(str)
	// 1,$(c volume-server list --comp) -> 1 $(c volume-server list --comp)
	pos, choice, str := stripPos(str)

	p.Position = pos
	p.PosChoice = choice
	p.Completion = str
	p.Help = help
	switch {
	case strings.HasPrefix(str, "<"):
		p.Type = Action
		p.Completion = strings.Trim(str, "<>")
	case strings.HasPrefix(str, "-"):
		p.Type = Option
	case strings.HasPrefix(str, "$("):
		p.Type = Command
	default:
		p.Type = String
	}
	return nil
}

// stripPos removes and saves a NUM,STRING, from the line.
func stripPos(str string) (int, string, string) {
	// if the first string up to a comma, looks like a number this might be something
	comma := strings.Index(str, ",")
	if comma < 0 {
		return 0, "", str
	}
	f, err := strconv.ParseUint(str[:comma], 10, 64)
	if err != nil {
		return 0, "", str
	}
	i := int(f)
	str = str[comma+1:]
	// do we have a 2nd comma?
	comma = strings.Index(str, ",")
	if comma < 0 {
		return i, "", str
	}
	// this must be a single string without spaces
	choice := str[:comma]
	if strings.Contains(choice, " ") {
		return i, "", str
	}
	str = str[comma+1:]
	return i, choice, str
}

// stripHelp check str for a [...] block at the end. If found that block is returned and removed from str, that new
// stripped help is then also returned.
func stripHelp(str string) (string, string) {
	if !strings.HasSuffix(str, "]") {
		return "", str
	}
	last := strings.LastIndex(str, "[")
	if last < 0 {
		return "", str
	}
	help := str[last:]
	str = str[:last]
	return help, str
}

// Cmd returns the "command" name from p. This is by definition the shortest key in p.
func (p Patterns) Cmd() string {
	cmd := ""
	for k := range p {
		if cmd == "" {
			cmd = k
		}
		if len(k) < len(cmd) {
			cmd = k
		}
	}
	return cmd
}

// Valid tells if p is valid. All keys must start with p.Cmd().
func (p Patterns) Valid() error {
	cmd := p.Cmd()
	for k := range p {
		if !strings.HasPrefix(k, cmd) {
			return fmt.Errorf("Key %q does not share the prefix: %q", k, cmd)
		}
	}
	return nil
}

func Tmpl(shell string) *template.Template {
	var err error
	tmpl := template.New(shell + ".go.tmpl") // .Funcs(ctx.FuncMap)
	tmpl, err = tmpl.ParseFS(tmplfs, shell+".go.tmpl")
	if err != nil {
		panic("cant find template: " + err.Error())
	}
	return tmpl
}

func quote(s string) string {
	if s == "" { // no need to quote empty strings...
		return s
	}
	return "'" + strings.TrimSpace(s) + "'"
}

func join(s []string) string {
	if len(s) == 0 {
		return ""
	}
	return " " + strings.Join(s, " ")
}
