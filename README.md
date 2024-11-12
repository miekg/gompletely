Re-implementation in Go of https://github.com/DannyBen/completely

Mostly to fix completion of positional arguments. See the README in the above repo for the YAML syntax (or
check the testdata directory).

The output closely matches 'completely', apart from the comments and the positional paramaters bit.
Positional parameters are only completed if they contain a command `$(...)` or are an action
`<file>`.

~~~ yaml
useradd:
- --root
- <noop>
- $(c group list --comp)
- <noop>
- <noop>
- $(c user lists --comp --contact)
- $(for m in 6 12 24 36 48; do ((m = m + 1)); echo $(date -d "$(date +%Y-%m-1) $m month" +%Y-%m-%d); done)
- $(echo E U S Z)
- $(c shell list --comp)
~~~

this `useradd` command, has 1 option (`--root`) and all other values are positional parameters with
a command used for the completion. The `<noop>` are noop completions to make postional counting
work.

For zsh completion

~~~ yaml
- "--root[help text"
~~~

is also supported. In the future we might even include ":message" suffix as well.
