Re-implementation in Go of https://github.com/DannyBen/completely

Mostly to fix completion of positional arguments. See the README in the above repo for the YAML syntax (or
check the testdata directory).

The output closely matches 'completely', apart from the comments and the positional paramaters bit.
Positional parameters are only completed if they contain a command `$(...)` or are an action
`<file>`.

Positional arguments are prefixed with an integer specifying which position they take.

~~~ yaml
useradd:
- '--root[blaa]'
- '2,$(c protogrp list --comp)'
- '5,$(c user lists --comp --contact)'
- '6,endate,$(for m in 6 12 24 36 48; do ((m = m + 1)); echo $(date -d "$(date +%Y-%m-1) $m month" +%Y-%m-%d); done)'
- '6,radboudid,$(echo E U S Z)'
- '6,shells,$(c shell list --comp)'
~~~

this `useradd` command, has 1 option (`--root`) and all other values are positional parameters with
a command used for the completion. The number in front of those is the position they can be in. If
the numbers are equal all of those positional argument can be expanded at that point. To distinguish
a unique (to this completion stanza) must be added.

The `[blaa]` is the optional help text that may be added. In the future we might even include
":message" suffix as well.
