##
# BASH Functions
##

test -n "$_ROOT" || exit 1

##
# Constants
##

COMMAND_SUCCESS=0
COMMAND_FAILURE=1
COMMAND_INVALID=2
COMMAND_CANNOT_EXECUTE=126
COMMAND_NOT_FOUND=127

# Log a message into a declared $LOGFILE
# log <message>
function log() {

	local _line
	if [ -n "$1" ]; then
		_line="$(date "+%Y-%m-%d %H:%M:%S") | $@"
	fi

	[ -n "$LOGFILE" ] || return
	if [ ! -e "$LOGFILE" ]; then
		local _dir="$(dirname "$LOGFILE")"
		[ -d "$_dir " ] || mkdir -p "$_dir" 2>/dev/null
		touch "$LOGFILE" 2>/dev/null || return
	fi
	if [ "$LOGVERBOSE" == "1" ] || [ "$LOGVERBOSE" == "true" ]; then
		echo "$_line" 2>&1
	fi
	if [ -n "$_line" ]; then
		echo "$_line" >>"$LOGFILE" 2>&1
	fi
}

# Display error message to stderr and exit
# throw [...message]
function throw() {
	if [ -n "$1" ]; then
		log "$@"
		echo "$@" >&2
	fi
	exit $COMMAND_CANNOT_EXECUTE
}

# Python pass
function pass() {
	return $COMMAND_SUCCESS
}

# split path separated by ':'
# splitpath <path>
function splitpath() {
	echo "$@" | tr ":" "\n"
}

##
# Functions
##

function function.copy() {
	if [ -n "$1" ] && [ -n "$2" ]; then
		if declare -F "$1" &>>/dev/null; then
			local func="$(declare -f "$1")"
			eval "${2}(${func#*\(}"
			return $COMMAND_SUCCESS
		fi
	fi
	return $COMMAND_FAILURE
}

# export_var <value> [name]
function export_var() {
	if [ -n "$2" ]; then
		declare -g "$2"="$1"
	else
		echo "$1"
	fi
	isset "$1"
}
# isset <var>
function isset() {
	[ -n "$1" ]
}

# mutes output
function @ {
	"$@" &>>/dev/null
}

#
# Includes a file
# require [...library]
function require() {
	[ -z "$1" ] && throw "require: Invalid argument count."
	local _script _filename _ext _arg

	until [ -z "$1" ]; do
		_arg="$1"
		_filename=""
		shift
		# file already has an extension
		if [[ "$_arg" =~ \.[a-zA-Z0-9]$ ]]; then
			_filename="$(script.find "$_arg")" || throw Cannot load library: "$_arg".
			require.load "$_filename"
			continue
		fi

		for _ext in $(splitpath "$REQUIRE_EXT"); do
			_script="${_arg}.${_ext}"
			if _filename="$(script.find "$_script")"; then
				require.load "$_filename"
				break
			fi
		done
		[ -n "$_filename" ] || throw Cannot load library: "$_arg".
	done
	return $COMMAND_SUCCESS
}

# require.load <filename>
function require.load() {
	[ -n "$1" ] || return $COMMAND_FAILURE
	[ -e "$1" ] || return $COMMAND_NOT_FOUND
	local _filename

	for _filename in "${REQUIRE_LOADED[@]}"; do
		[ "$_filename" == "$1" ] && return $COMMAND_SUCCESS
	done
	REQUIRE_LOADED[${#REQUIRE_LOADED[@]}]="$1"

	pushd "$(dirname "$1")" &>>/dev/null
	source "$1"
	popd &>>/dev/null
	return $COMMAND_SUCCESS
}

# create random UUID
function UUID.create() {
	cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
}

#
# Get uptime in seconds
# https://gist.github.com/OndroNR/0a36f97cd612b75fbf92f22cf72851a3
#
function uptime.int() {
	if [ -e /proc/uptime ]; then
		cat /proc/uptime | awk '{printf "%0.f", $1}'
	else
		set +e

		if sysctl kern.boottime &>/dev/null; then
			local kern_boottime=$(sysctl kern.boottime 2>/dev/null | sed "s/.* sec\ =\ //" | sed "s/,.*//")
			local time_now=$(date +%s)
			local uptime=$((time_now - kern_boottime))
			echo $uptime
		else echo "-1"; fi
		set -e
	fi
}

# This is a general-purpose function to ask Yes/No questions in Bash, either
# with or without a default answer. It keeps repeating the question until it
# gets a valid answer.
# Source:
# https://gist.github.com/davejamesmiller/1965569
function ask() {

	[ "${OVERRIDE}" = "yes" ] && return 0

	# https://djm.me/ask
	local prompt default reply

	while true; do

		if [ "${2:-}" = "Y" ]; then
			prompt="Y/n"
			default=Y
		elif [ "${2:-}" = "N" ]; then
			prompt="y/N"
			default=N
		else
			prompt="y/n"
			default=
		fi

		# Ask the question (not using "read -p" as it uses stderr not stdout)
		echo -n "$1 [$prompt] "

		# Read the answer (use /dev/tty in case stdin is redirected from somewhere else)
		read reply

		# Default?
		if [ -z "$reply" ]; then
			reply=$default
		fi

		# Check if the reply is valid
		case "$reply" in
		Y* | y*) return 0 ;;
		N* | n*) return 1 ;;
		esac

	done
}

# Emulate Batch pause
# pause [timeout]
function pause() {
	local input
	local timeout=0
	local ret=0
	[ -n "$1" ] && echo "$1" | grep -Eq "^[0-9]+$" && timeout=$1
	if [ "x$timeout" != "x0" ]; then
		until [ $timeout -eq 0 ]; do
			printf "\rPlease hit a key to continue (Wait %d seconds or hit Esc to cancel) ..." $timeout
			((timeout = $timeout - 1))
			read -t1 -n1 -s input
			ret=$?
			[ $ret -lt 2 ] && break
		done
		printf "\n"

	else
		echo "Please hit a key to continue (Esc to cancel) ..."
		read -n1 -s input
		ret=$?
	fi
	[ $ret -gt 0 ] && ret=1
	case $input in
	$'\e') ret=1 ;;
	esac
	return $ret
}
