# bash completion script for AddUser
# shellcheck disable=SC2207

# logmaar() {
#     cat - >> /tmp/logmaar
# }

_dbcomplete_protogroup() {
	local cur opts
	opts="$(/usr/local/bin/dblist -i --protogroups)"
	cur="${COMP_WORDS[COMP_CWORD]}"
	COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
}

_dbcomplete_ictcontact() {
	local cur opts
	opts="$(/usr/local/bin/dblist -i --ict-contacts)"
	cur="${COMP_WORDS[COMP_CWORD]}"
	COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
}

_complete_checkdate() {
	local dates=""
	# generate possible check dates for the login
	for m in 6 12 24 36 48; do
		((m = m + 1))
		dates="$dates $(date -d "$(date +%Y-%m-1) $m month" +%Y-%m-%d)"
	done
	COMPREPLY=($(compgen -W "$dates" -- ${cur}))
}

_dbcomplete_shell() {
	local cur opts
	opts="$(/usr/local/bin/dblist -i --shells)"
	cur="${COMP_WORDS[COMP_CWORD]}"
	COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
}

_userdb_adduser_completion() {
	COMPREPLY=()
	local cur

	# onderstaande functie is gedefinieerd in  /usr/share/bash-completion/bash_completion
	# en zorgt ervoor dat de -n opgegeven karakters niet als scheidingsteken voor
	# de argumenten worden gebruikt.
	_get_comp_words_by_ref -n '@><=;|&(:' cur prev
	if [[ ${cur} == -* ]]; then
		# current word is starting with a dash: filling COMPREPLY with all available flags
		COMPREPLY=($(compgen -W "--reuse_uid --no-reuse_uid -man --help --root --homevolume --no-homevolume --mailcontact --maildhzpasswordto --no-mailcontact" -- $cur))
		return
	else
		# determine positional argument count
		posargcount=$COMP_CWORD
		for i in "${COMP_WORDS[@]}"; do
			if [[ ${i} == -* ]]; then
				((posargcount = posargcount - 1))
			fi
		done

		case "$posargcount" in
		2)
			# protogroup
			_dbcomplete_protogroup
			return
			;;
		5)
			# ict contact
			_dbcomplete_ictcontact
			return
			;;
		6)
			# U- of S-nummer
			return
			;;
		7)
			# check date
			_complete_checkdate
			return
			;;
		8)
			# shell
			_dbcomplete_shell
			return
			;;
		esac
	fi
}

complete -F _userdb_adduser_completion AddUser
