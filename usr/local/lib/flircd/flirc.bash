# FLIRC Utilities

flirc_util="$(script.find flirc_util)"
[ -z "$flirc_util" ] && throw "Cannot find flirc_util"

function flirc() {
    [ -z "$1" ] && return 1
    flirc_util $@
}

# to be executed in the main thread
# detect who is connected to the display and which display
flirc.init.display() {
    local _user _disp

    if [ -z "$flirc_user" ]; then
        for info in $(who); do
            # echo $info
            if [ -z "$_user" ]; then
                _user="$info"
                continue
            fi
            if [ "$info" != "${info#(}" ]; then
                info="${info#(}"
                info="${info%)}"
                if [ -z "${info#:[0-9]}" ]; then
                    _disp="$info"
                    break
                else
                    _user=""
                fi
            fi
        done
        if [ -z "$_user" ] || [ -z "$_disp" ]; then
            return $COMMAND_FAILURE
        fi
        export flirc_user="$_user"
        export flirc_display="$_disp"

        if [ "$DISPLAY" != "$flirc_display" ]; then
            export DISPLAY="$flirc_display"
        fi
    fi

    return $COMMAND_SUCCESS
}

# execute command with display unlocked
flirc.sudo() {
    [ -z "$1" ] && return $COMMAND_FAILURE
    sudo -u $flirc_user $@
}

# Checks if device is connected
function flirc.connected() {
    @ flirc.devices
}

# Finds flirc devices and display their path
flirc.devices() {
    local devices result dev
    find /dev/input/by-id | grep "flirc" 2>/dev/null
}

# Finds flirc xinput ids
flirc.ids() {
    local ids i
    for i in $(flirc.sudo xinput list | grep "flirc"); do

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
    if echo $status | grep "SKU" &>>/dev/null; then
        text="connected"
        result=0
    elif echo $status | grep "Bootloader" &>>/dev/null; then
        text="bootloader"
        result=2
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

    # other shell detections are to be added here
    local _processes="plasmashell gnome-shell lxqt-session xfce4-session gnome-session-binary" _ids _id _result=$COMMAND_FAILURE _mode=disable

    if [ "$1" == "false" ]; then
        _mode=enable
    fi

    # instantiate only once
    if ! instanceof desktops ProcessList; then
        new ProcessList desktops $_processes
    fi
    if flirc.connected && desktops.running; then
        flirc.init.display
        _ids="$(flirc.ids)"
        if [ -n "$_ids" ]; then
            log "flirc device connected and desktop is running"

            for _id in $_ids; do
                if flirc.sudo xinput ${_mode} $_id; then
                    _result=$COMMAND_SUCCESS
                    log "flirc device xinput[id=${_id}] has been ${_mode}d"
                else
                    log "ERR: cannot ${_mode} device xinput[id=${_id}]"
                fi

            done
        else
            log "ERR: cannot detect flirc xinput device ids"
        fi
    fi
    return $_result
}
