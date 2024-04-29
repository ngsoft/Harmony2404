#
# Manages FLIRC
#

>/dev/null 2>&1 which flirc_util || throw.log Cannot find flirc_util

flirc(){
    [ -z "$1" ] && return 1
    flirc_util $@
}

flirc.connected(){
    >/dev/null 2>&1 flirc.devices
}

# Finds flirc devices and display their path
flirc.devices(){
    local devices result dev
    find /dev/input/by-id | grep "flirc" 2> /dev/null
}
# Finds flirc xinput ids
flirc.xinput.ids(){
    local ids i
    for i in $(xinput list | grep "flirc"); do
        if [[ "$i" =~ id\=([0-9]+) ]]; then 
            [ -z "$ids" ] || ids+=" "
            ids+=${BASH_REMATCH[1]}
        fi
    done
    [ -z "$ids" ] && return 1
    echo $ids
}

# Get Flirc current status 
# return $? = 1 : disconnected
# return $? = 2 : Bootloader
# return $? = 0 : connected
flirc.status(){
    local status text result
    status=$(2>/dev/null flirc version)
    text=disconnected
    result=1
    echo $status | grep "SKU" >/dev/null 2>&1  && text="connected"; result=0
    echo $status | grep "Bootloader" > /dev/null 2>&1 && text="bootloader"; result=2
    echo $text
    return $result
}

# Fix FLIRC in bootloader mode
flirc.fixdfu(){
    flirc dfu leave | grep "FW Detected" > /dev/null 2>&1
}

# Set configuration options flirc
flirc.set(){
    [ -n "$1" ] || return 1
    [ -n "$2" ] || return 1
    local option=$1
    local value=$2
    local valid=(sleep_detect seq_modifiers noise_canceler profiles interkey_delay)
    if array.contains $option ${valid[@]}; then
        if [ "$option" == "interkey_delay" ]; then
            [[ $value =~ ^[0-6]$ ]] || return 1
            >/dev/null 2>&1 flirc $option $value
            return $?
        elif [[ $value =~ (disable|enable) ]]; then
            >/dev/null 2>&1 flirc $option $value
            return $?
        fi
    fi
    return 1
}