#
# Helper Library to design bash scripts
#
if [ -z "$HARMONY_ENV_LOADED" ]; then
    echo Harmony Environment is not loaded >&2
    exit 1
fi


file.find(){
    is.empty $1 && return 1
    local dir
    for dir in $(echo $HARMONY_INCLUDE_PATH | tr ":" "\n"); do
        if [ -e "$dir/$1" ]; then 
            echo "$dir/$1"
            return 0
        fi
    done
    return 1
}

#check if array contains value
# array.contains <needle> <${haystack[@]}>
array.contains() {
    [ -n "$1" ] || return 1
    local needle=$1
    shift
    for e in "$@"; do [[ "$e" == "$needle" ]] && return 0; done
    return 1
}

# Array utils

# Exports Array to file
# array.export <varname> <filename>
array.export(){
    is.empty $1 && return 1
    is.empty $2 && return 1
    local var named="_exported" file dir key value line
    file="$2"
    dir="$(dirname "$2")"
    [ ! -e "$dir" ] && cmd mkdir -p "$dir"
    [ ! -d "$dir" ] && return 1
    
    declare -n var=$1
    if [ "${#var[@]}" != "0" ]; then
        echo "$named=()">"$file"
        for key in "${!var[@]}"; do
            value="${var[key]}"
            line="$named"
            line+="[$key]="
            line+='"'
            line+="$value"
            line+='"'
            echo "$line">>"$file"
        done
        return 0
    fi
    return 1
}

# Imports Array from file exported using array.export
# array.export <varname> <filename>
array.import(){
    is.empty $1 && return 1
    is.empty $2 && return 1
    local _exported=() file var="$1" key value
    file="$2"
    [ -e "$file" ] || return 1
    source "$file"
    [ "${#_exported[@]}" == "0" ] && return 1
    declare -g -A "$var"
    declare -n var="$var"
    for key in "${!_exported[@]}"; do
        value="${_exported[$key]}"
        var[$key]="$value"
    done
    return 0
}



# parse.configfile <file> <outvariable>
parse.configfile(){
    local src=$1
    local out=$2
    local params
    [ -n "$2" ] || return 1
    declare -A -g $out
    declare -n params=$out
    for file in $src; do
        [ -f $file ] || continue
        while read aline; do
            aline=${aline//\#*/}
            [[ -z $aline ]] && continue
            read var value <<<$aline
            [[ -z $var ]] && continue
            params[$var]=${value//[\=\ ]/''}
        done <$file
    done
}

# Executes given command 
# and hides output
cmd(){
    [ -n "$1" ] || return 1
    local _cmd
    while [ -n "$1" ]; do
        [ -z "$_cmd" ] || _cmd+=' '
        if [[ "$1" =~ [\ ]+ ]]; then _cmd+='"'$1'"'; else _cmd+=$1; fi
        shift
    done
    > /dev/null 2>&1 eval $_cmd
}

# Fake process
noop(){
    return 0
}

UUID() {
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1
}

function.copy() {
    if [ -n "$1" ] && [ -n "$2" ]; then
        declare -F "$1" > /dev/null || return 1
        local func="$(declare -f "$1")"
        eval "${2}(${func#*\(}"
        return 0
    fi
    return 1
}

new() {
    [ -n "$1" ] || throw Invalid constructor
    [ -n "$2" ] || throw Invalid instance
    local this base dir method filename script fn destr constructor retval=0
    constructor=; base=$1; this=$2; script=$base.class.sh; shift 2;
    filename=$(file.find "$script")
    [ -n "$filename" ] || throw Cannot find class $base
    fn='this.prop(){ if [ "$2" == "=" ]; then _this[$1]=$3; elif is.empty $1; then return 1; elif ! is.empty $2; then _this[$1]=$2; else echo ${_this[$1]}; fi }'
    fn=${fn//this/$this}
    destr='this.destroy(){ local method var; unset -v _this; for method in $(compgen -A function this); do unset -f ${method}; done; for var in $(compgen -A variable this); do unset -v ${var}; done }'
    destr=${destr//this/$this}
    source <(sed "s/this/$this/g" $filename)
    if cmd declare -F $base; then
        declare -g -A _${this}
        if ! cmd declare -F ${this}.prop; then eval $fn; fi
        if ! cmd declare -F ${this}.destroy; then eval $destr; fi
        cmd pushd "$(dirname "$filename")"
            $base
            constructor=
            for method in ${this}_constructor ${this}.constructor;do 
                cmd declare -F $method && constructor=$method
            done
            if [ -n "$constructor" ]; then 
                $constructor $@
                retval=$?
                unset -f $constructor
            fi
        cmd popd
    else throw "Invalid class declaration $base."; fi
    return $retval
}


# Log a message into a declared $logfile
# log <message>
log(){
    [ -n "$logfile" ] || return
    if [ ! -e "$logfile" ]; then
        local dir=$(dirname "$logfile")
        [ -d "$dir " ] || mkdir -p $dir 2>/dev/null
        touch $logfile 2>/dev/null
        [ $? -gt 0 ] && return
    fi
    if [ -n "$1" ]; then 
        echo "$(date "+%Y-%m-%d %H:%M") | $@" >> $logfile 2>&1; 
    fi
}

# Display error message to stderr and exit
throw(){
    [ -z "$1" ] || echo "$@" >&2
    exit 1
}
throw.log(){
    log $@
    throw $@
}


# Uses Ubuntu Notification system to notify the user
# and logs the message
# <BODY> [ICON]
notify(){
    [ -n "$1" ] || return 1
    local icon body
    icon=dialog-information
    if [ $# -eq 2 ]; then icon=$2; body=$1
    elif [ $# -eq 1 ]; then body=$1
    else return 1; fi
    notify-send --urgency=normal --icon=$icon "$body"
}
# Notify and logs the message
notify.log(){
    [ -n "$1" ] || return 1
    log $1
    notify "$1" $2
    return $?
}

# Display error message to stderr, notify and exit
notify.throw(){
    notify "$1" "$2"
    throw "$1"
}

# Require a file
require(){
    [ -z "$1" ] && throw require invalid argument
    local script filename
    until [ -z "$1" ]; do
        filename=
        if [[ $1 =~ \.[a-zA-Z0-9]{1,3}$ ]]; then script=$1; else script=$1.sh; fi
        filename=$(file.find "$script")
        [ -n "$filename" ] || throw Cannot find library $1
        source $filename
        shift
    done
    return 0
}

# Send signal to process
# signal.send <process> ...<signal>
signal.send(){
    [ -n "$1" ] || return 1
    [ -n "$2" ] || return 1
    local pid process
    process=$1
    shift
    until [ -z "$1" ]; do
        for pid in $(pidof -x $process); do
            kill -s $1 $pid
        done
        shift
    done
}



# declare an event
# on <type> <action>
on(){
    declare -g -A stack
    [ -n "$1" ] || return 1
    [ -n "$2" ] || return 1
    [[ $1 =~ ^[a-zA-Z][a-zA-Z0-9_]+$ ]] || return 1
    local type var
    type=$1
    shift
    if [ -z "${stack[$type]}" ]; then
        var=emit_$type
        declare -g -a $var
        stack[$type]=$var
    fi
    var=${stack[$type]}
    eval $var+="(\""$@"\")"
    return 0
}



# Trigger an event
# trigger ...<type>
trigger(){
    declare -g -A stack
    [ -n "$1" ] || return 1
    local type fn var result tmp
    until [ -z "$1" ]; do
        result=1
        type="$1"
        #legacy create a on.type()
        fn="on.$type"
        if cmd declare -F $fn; then result=0; $fn; fi 
        #use the stack
        var=${stack[$type]}
        if [ -n "$var" ]; then
            declare -n tmp=$var
            for fn in "${tmp[@]}"; do 
                if cmd declare -F $fn; then cmd $fn; result=$?; fi
            done
        fi
        if [ $result -gt 0 ]; then log Event $type got an error; fi
        shift
    done
    return 0
}

#
# Get uptime in seconds
# https://gist.github.com/OndroNR/0a36f97cd612b75fbf92f22cf72851a3
#
function uptime.int()
{
  if [ -e /proc/uptime ] ; then
    echo `cat /proc/uptime | awk '{printf "%0.f", $1}'`
  else
    set +e
    sysctl kern.boottime &> /dev/null
    if [ $? -eq 0 ] ; then
      local kern_boottime=`sysctl kern.boottime 2> /dev/null | sed "s/.* sec\ =\ //" | sed "s/,.*//"`
      local time_now=`date +%s`
      local uptime=$(($time_now - $kern_boottime))
      echo $uptime
    else
      echo "-1"
    fi
    set -e
  fi
}

# This is a general-purpose function to ask Yes/No questions in Bash, either
# with or without a default answer. It keeps repeating the question until it
# gets a valid answer.
# Source:
# https://gist.github.com/davejamesmiller/1965569
ask() {

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
            Y*|y*) return 0 ;;
            N*|n*) return 1 ;;
        esac

    done
}

# Emulate Batch pause
# pause [timeout]
pause(){
    local input
    local timeout=0
    local ret=0
    [ -n "$1" ] && echo "$1" | grep -Eq ^[0-9]+$  && timeout=$1
    if [ "x$timeout" != "x0" ]; then
        until [ $timeout -eq 0 ]; do
            printf "\rPlease hit a key to continue (Wait $timeout seconds or hit Esc to cancel) ..."
            ((timeout=$timeout-1))
            read -t1 -n1 -s input
            ret=$?
            [ $ret -lt 2 ] && break;
        done 
        printf "\n"
        
    else 
        echo "Please hit a key to continue (Esc to cancel) ..."
        read -n1 -s input
        ret=$?
    fi
    [ $ret -gt 0 ] && ret=1
    case $input in
        $'\e' ) ret=1;;
    esac
    return $ret
}

# Get all the pids
process.pids(){
    [[ -z $1 ]] && return 1
    pgrep -f "$1"
}

# Count precess running
process.count(){
    [[ -z $1 ]] && return 1
    pgrep -c -f "$1"
}


#get the process pid
process.pid(){
    [[ -z $1 ]] && return 1
    pgrep -n -f "$1"
}

# Check if a process is running
process.running(){
    [ -z "$1" ] && return 1
    cmd pidof -x $1
}



# Search for process pids (if any) and kills them
process.kill(){
    local pids _pid
    [[ -z $1 ]] && return 1
    process.running $1 || return 0
    cmd killall $1
}

# process.wait <process> [timeout=10s]
process.wait(){
    local timeout=10 cnt=0 process
    is.empty $1 && return 1
    process=$1
    is.int.unsigned $2 && timeout=$2
    ((timeout=$timeout*5))
    while true; do
        sleep .2
        ((cnt+=1))
        cmd pgrep $process && return 0
        if [[ $cnt -gt $timeout ]]; then return 1; fi
    done
    return 0;
}


#is methods
is.bool(){
    [[ $1 =~ (true|false) ]]
}

is.true(){
    [ "$1" == "true" ]
}

is.int(){
    [[ $1 =~ ^[\+\-]?[0-9]+$ ]]
}
is.int.unsigned(){
    [[ $1 =~ ^[0-9]+$ ]]
}
is.double(){
    [[ $1 =~ ^[\+\-]?[0-9]*\.[0-9]+$ ]]
}
is.double.unsigned(){
    [[ $1 =~ ^[0-9]*\.[0-9]+$ ]]
}
is.empty(){
    [ -z "$1" ]
}
is.function(){
    cmd declare -F $1
}


HARMONY_UTILS_LOADED=1