# FLIRC Utilities

flirc_util="$(script.find flirc_util)"
[ -z "$flirc_util" ] && throw "Cannot find flirc_util"

function flirc() {
    [ -z "$1" ] && return 1
    flirc_util $@
}

# Checks if device is connected
function flirc.connected() {
    if @ flirc.devices; then
        if [ -z "$flirc_model" ]; then
            flirc_model=1
            if flirc version | grep SKU | grep 2.0 &>>/dev/null; then
                flirc_model=2
            fi
            export flirc_model
        fi
        return $COMMAND_SUCCESS
    fi
    export flirc_model=0
    return $COMMAND_FAILURE
}

# Finds flirc devices and display their path
flirc.devices() {
    local devices result dev
    find /dev/input/by-id | grep "flirc" 2>/dev/null
}

# Finds flirc xinput ids
flirc.ids() {
    local ids i
    for i in $(sudo.execute xinput list | grep "flirc"); do

        if [[ "$i" =~ id\=([0-9]+) ]]; then
            [ -z "$ids" ] || ids+=" "
            ids+=${BASH_REMATCH[1]}
        fi
    done
    [ -z "$ids" ] && return $COMMAND_FAILURE
    echo $ids
}

# Get Flirc device current status
# return $? = 1 : disconnected
# return $? = 2 : Bootloader
# return $? = 0 : connected
flirc.status() {
    local status text result
    status=$(flirc 2>/dev/null version)
    text="disconnected"
    result=1

    if echo $status | grep "Bootloader" &>>/dev/null; then
        text="bootloader"
        result=2
    elif echo $status | grep "SKU" &>>/dev/null; then
        text="connected"
        result=0
    fi
    echo $text
    return $result
}

# Fix FLIRC in bootloader mode
flirc.dfu() {
    flirc dfu leave | grep "FW Detected" &>>/dev/null
}

# Set configuration options flirc
flirc.set() {
    [ -n "$1" ] || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE
    local option=$1
    local value=$2
    local valid=(sleep_detect seq_modifiers noise_canceler profiles interkey_delay)
    if array.contains $option ${valid[@]}; then
        if [ "$option" == "interkey_delay" ]; then
            if [ "$flirc_model" == "2" ]; then
                return $COMMAND_FAILURE
            fi
            [[ $value =~ ^[0-6]$ ]] || return 1
            flirc $option $value &>>/dev/null
            return $?
        elif [[ $value =~ (disable|enable) ]]; then
            flirc $option $value &>>/dev/null
            return $?
        fi
    fi
    return $COMMAND_FAILURE
}

# Detect desktop shell running and lock xinput events for the device (it is a keyboard first)
# <kde> <gnome> <lxde> <xfce> <unity>
# flirc.capture [disable:true|false]
flirc.capture() {

    local mode=disable result=$COMMAND_FAILURE ids id
    if is false "$1" || is off "$1"; then
        mode=enable
    fi

    if flirc.connected && session.detect; then
        if ids="$(flirc.ids)"; then
            log "flirc device connected and desktop is running"
            for id in $ids; do
                if @ sudo.execute xinput $mode $id; then
                    result=$COMMAND_SUCCESS
                    log "flirc xinput device ${id} has been ${mode}d"
                    continue
                fi
                log "ERR: cannot ${mode} xinput device ${id}"
            done
        fi
    fi
    return $result
}
