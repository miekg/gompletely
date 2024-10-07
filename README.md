Re-implementation in Go of https://github.com/DannyBen/completely

Mostly to fixes completion of positional arguments. See the README in the above for the YAML syntax (or
check the testdata directory).

Allow the following extra syntax:

~~~ yaml
command:
- <noop>
~~~

Where `<noop>` is a null action that is used to count positional argument and have the completion
work for those too.
