package main

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Patterns is a whole file of completions for a command.
type Patterns map[string][]Pattern

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
type Type int

const (
	None Type = iota
	Action
	Option
	Command
	String
)

// Pattern is a completion we read from the yaml. It is altered and made suitable for completion
// generation by Bash/Zsh/... etc.
type Pattern struct {
	CompType Type
	CompGen  string
	Case     string
}

func (c *Pattern) UnmarshalYAML(node *yaml.Node) error {
	str := ""
	err := node.Decode(&str)
	if err != nil {
		return err
	}

	c.CompGen = str
	switch {
	case strings.HasPrefix(str, "<"):
		c.CompType = Action
	case strings.HasPrefix(str, "-"):
		c.CompType = Option
	case strings.HasPrefix(str, "$("):
		c.CompType = Command
	default:
		c.CompType = String
	}
	return nil
}
