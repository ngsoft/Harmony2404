# process.wait <process> [timeout=10s]
process.wait() {
    local timeout=10 cnt=0 process
    is.empty $1 && return 1
    process=$1
    is.unsigned $2 && timeout=$2
    ((timeout = $timeout * 5))
    while true; do
        sleep .2
        ((cnt += 1))
        @ pgrep $process && return $COMMAND_SUCCESS
        if [[ $cnt -gt $timeout ]]; then return $COMMAND_FAILURE; fi
    done
    return $COMMAND_SUCCESS
}

# Search for process pids (if any) and kills them
process.kill() {
    local pids _pid
    [[ -z $1 ]] && return $COMMAND_FAILURE
    process.running $1 || return $COMMAND_SUCCESS
    @ killall $1
}

# Check if a process is running
process.running() {
    [ -z "$1" ] && return $COMMAND_FAILURE
    @ pidof -x $1
}

#get the process pid
process.pid() {
    [[ -z $1 ]] && return $COMMAND_FAILURE
    pgrep -n -f "$1"
}

# Count precess running
process.count() {
    [[ -z $1 ]] && return $COMMAND_FAILURE
    pgrep -c -f "$1"
}
