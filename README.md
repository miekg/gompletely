Re-implementation in Go of https://github.com/DannyBen/completely

> Note: Hacked together (esp the Zsh completion), works for me, probably not for you.

This generates Zsh and Bash completions from a YAML definition.

This mostly fixes completion of positional arguments. See the README in the above repo for the YAML
syntax (or check the testdata directory). Although for Bash I feel there needs to be a better way -
seems to work just good enough. Note this YAML is very much a hacked up version from the original,
the value encode also much more - in fact is almost not really YAML anymore.

The output closely matches 'completely' (for Bash), apart from the comments and the positional
parameters bit. Positional parameters are only completed if they contain a command `$(...)` or are
an action `<file>`. The action are the action as defined for Bash's `compgen` and are converted to
Zsh actions. (This might reverse in the future).

Positional arguments are prefixed with an integer specifying which position they take, that must be
in sequence, starting with 1. Zsh completion silently fails if this isn't the case.

If a positional argument does not have a completion you can let Zsh say what you need to "complete"
(=to type in) there with `<int>,<message>,` so that on \<TAB\>, \<message\> will be displayed.

In brackets any help message may be put: `[help message]`, the Zsh completion will show that. This
help message must be the last things on the line in the YAML. Backticks in the help text are
replaced with U+FF40 characters (Fullwidth Grave Accent).

If you have a subcommand, which is a kind of a positional argument that has several choices in that
(same) position use `S,subcommand,`.

The Zsh part of this is under active development, I don't use bash.

The help text is assumed to be in markdown, which is stripped to prevent weird character showing up
in the shell completion scripts.

~~~ yaml
useradd:
- '--root[help message]'
- '1,$(c protogrp list --comp)'
- '2,message,
- '3,endate,$(for m in 6 12 24 36 48; do ((m = m + 1)); echo $(date -d "$(date +%Y-%m-1) $m month" +%Y-%m-%d); done)'
- '4,shells,$(c shell list --comp)'
~~~

## Development

### Zsh

Unload and load test completion:
~~~ sh
unfunction _AddVolume; autoload -U _AddVolume
./gompletely -s zsh < testdata/AddVolume.yml > _AddVolume
source _AddVolume
~~~

### TODO

* Tests
* Fix the coding mess that is zsh.go
* Tests against previously generated *.bash and \_zsh completion files.
