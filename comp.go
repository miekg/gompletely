package main

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Definition is a whole file of completion for a command.
type Definition map[string][]Completion

// Completion holds a single line of the completion. One of these is non-nil and should be used.
type Completion struct {
	Action
	Option
	Command
	String
}

// An Action is a bash shell action to e.g. complete files. These are denoted in the yaml as "<action>".
// There is one special action "<none>" which is used for arguments that don't have any completion.
type Action string

// An Option is the option that can be completed for a command, i.e. "-f" or "--force". These are recognized because
// they all start with a minus,
type Option string

// Command is a shell command to be executed, these are recognized because they use the command substituion
// syntax "$(echo hello)"
type Command string

// String is a single string that is the completion.
type String string

func (c *Completion) UnmarshalYAML(node *yaml.Node) error {
	str := "" // they are all strings
	err := node.Decode(&str)
	if err != nil {
		return err
	}
	switch {
	case strings.HasPrefix(str, "<"):
		c.Action = Action(str)
	case strings.HasPrefix(str, "-"):
		c.Option = Option(str)
	case strings.HasPrefix(str, "$("):
		c.Command = Command(str)
	default:
		c.String = String(str)
	}

	return nil
}
