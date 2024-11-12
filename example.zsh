
#compdef c
# zsh completion for c

_c_user_(u)_list_(ls)() {
    _arguments -C \
        "--no-header[Do not output a header when listing when listing all users.]" \
        '(-f --full)'{-f,--full}"[Display full output of the user.]" \
        '(-T --terse)'{-T,--terse}"[Output the user in a terse format.]" \
        '(-r --removed)'{-r,--removed}"[Display a removed user account. This requires access to /var/log/removed-users or present the JSON on standard input.]" \
        "--contact[When listing all users show only those that are a contact for another user.]" \

}
