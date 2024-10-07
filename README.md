Re-implementation in Go of https://github.com/DannyBen/completely

Mostly to fix completion of positional arguments. See the README in the above repo for the YAML syntax (or
check the testdata directory).

The output closely matches 'completely', apart from the comments and the positional paramaters bit.
Positional parameters are only completed if they contain a command `$(...)` or are an action
`<file>`. Options and lists of strings are not included. For subcommands this probably doesn't work.
