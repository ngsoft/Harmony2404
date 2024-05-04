##
# Array utils
##

function array() {
    pass
}

# array.create <name>
function array.exists() {
    [ -n "$1" ] || return $COMMAND_FAILURE
    [[ "$(declare -p "$1" 2>/dev/null)" =~ declare[^-]-[aA] ]]
}

# create new empty array
# array.create <name>
function array.create() {
    [ -n "$1" ] || return $COMMAND_FAILURE
    eval "declare -g -A ${1}=()"
}

# create new indexed empty array
# array.create.list <name> [...entries]
function array.create.list() {
    [ -n "$1" ] || return $COMMAND_FAILURE
    local _name=${1}
    shift
    eval "declare -g -a ${_name}=()"
    if [ -n "$1" ]; then
        array.push ${_name} $@
    fi
}

# array.count <name>
function array.count() {
    array.exists "$1" || return $COMMAND_FAILURE
    local _array
    declare -n _array="$1"
    echo "${#_array[@]}"
    return $COMMAND_SUCCESS
}

# Gets the first key of an array
# array.last <name>
function array.first() {
    array.exists "$1" || return $COMMAND_FAILURE
    local _name _key _prev
    declare -n _name="$1"
    for _key in "${!_name[@]}"; do
        pass
    done
    if [ -n "$_key" ]; then
        echo "$_key"
        return $COMMAND_SUCCESS
    fi

    return $COMMAND_FAILURE
}

# Gets the last numeric key of an array
# array.last <name>
function array.last() {
    array.exists "$1" || return $COMMAND_FAILURE
    local _name _key
    declare -n _name="$1"
    for _key in "${!_name[@]}"; do
        echo "$_key"
        return $COMMAND_SUCCESS
    done
    return $COMMAND_FAILURE
}

# dumps array
# array.export <name>
function array.dump() {
    array.exists "$1" || return $COMMAND_FAILURE

    local _dump

    _dump=$(declare -p "$1")

    if [[ "$_dump" =~ \((.*)\)$ ]]; then
        echo "${BASH_REMATCH[1]}"
        return $COMMAND_SUCCESS
    fi

    return $COMMAND_FAILURE

}

# array.push <name> [...values]
function array.push() {

    array.exists "$1" || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE

    local _old _new _key _index
    declare -n _old="$1"
    _new=()

    for _key in "${!_old[@]}"; do
        if is.int "$_key"; then
            _new[${#_new[@]}]="${_old[$_key]}"
            unset "_old[$_key]"
        fi
    done
    _index="${#_new[@]}"
    until [ -z "$2" ]; do
        _new[${#_new[@]}]="$2"
        ((_index += 1))
        shift
    done
    for _key in "${!_new[@]}"; do
        _old[$_key]="${_new[$_key]}"
    done

    return $COMMAND_SUCCESS
}

# array.unshift <name> [...values]
array.unshift() {

    array.exists "$1" || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE
    local _new _old _key _tmp
    declare -n _old="$1"
    _new=()
    _tmp=()

    # adds new values
    until [ -z "$2" ]; do
        _new[${#_new[@]}]="$2"
        shift
    done
    # adds old values _tmp
    for _key in "${!_old[@]}"; do
        if is.int "$_key"; then
            _tmp[${#_tmp[@]}]="${_old["$_key"]}"
            unset "_old["$_key"]"
        fi
    done
    # reverse
    for _key in $(echo "${!_tmp[@]}" | rev); do
        _new[${#_new[@]}]="${_tmp[$_key]}"
    done

    #merge
    for _key in "${!_new[@]}"; do
        _old[$_key]="${_new[$_key]}"
    done

}

# array.pop <name> [varname]
function array.pop() {
    array.exists "$1" || return $COMMAND_FAILURE
    local _name _key
    declare -n _name="$1"
    for _key in "${!_name[@]}"; do
        if [ -n "$2" ]; then
            declare -g "$2"="${_name["$_key"]}"
        fi
        unset "_name[$_key]"
        return $COMMAND_SUCCESS
    done
    return $COMMAND_FAILURE
}

#check if array contains value
# array.contains <name> <needle>
function array.contains() {

    if [ -z "$1" ] || [ -z "$2" ]; then
        return $COMMAND_FAILURE
    fi
    array.exists "$1" || return $COMMAND_FAILURE
    local _needle _haystack
    declare -n _haystack="$1"

    for _needle in "${_haystack[@]}"; do
        [ "$_needle" == "$2" ] && return $COMMAND_SUCCESS

    done
    return $COMMAND_FAILURE
}

# search key for value
# array.contains <name> <needle> [key]
function array.search() {

    if [ -z "$1" ] || [ -z "$2" ]; then
        return $COMMAND_FAILURE
    fi
    array.exists "$1" || return $COMMAND_FAILURE
    local _needle _haystack _key
    declare -n _haystack="$1"

    for _key in "${!_haystack[@]}"; do
        _needle="${_haystack[$_key]}"
        if [ "$_needle" == "$2" ]; then
            if [ -n "$3" ]; then
                declare -g "$3"="$_key"
            else
                echo "$_key"
            fi
            return $COMMAND_SUCCESS
        fi
    done
    return $COMMAND_FAILURE
}

# array.shift <name> [varname]
function array.shift() {
    array.exists "$1" || return $COMMAND_FAILURE
    local _name _key
    declare -n _name="$1"
    for _key in "${!_name[@]}"; do
        pass
    done

    if [ -n "$_key" ]; then
        if [ -n "$2" ]; then
            declare -g $2="${_name["$_key"]}"
        fi
        unset "_name[$_key]"
        return $COMMAND_SUCCESS
    fi

    return $COMMAND_FAILURE
}

# array.filter <name> <new_name> <function_name>
function array.filter() {
    declare -F "$3" &>>/dev/null || return $COMMAND_FAILURE
    array.exists "$1" || return $COMMAND_FAILURE
    [ "$1" == "$2" ] && return $COMMAND_FAILURE
    array.create "$2" || return $COMMAND_FAILURE

    local _key _result _array

    declare -n _array="$1"
    declare -n _result="$2"

    for _key in "${!_array[@]}"; do
        if $3 "${_array[$_key]}" "$_key"; then
            _result[$_key]="${_array[$_key]}"
        fi
    done
    return $COMMAND_SUCCESS
}

# array.get <name> <offset> [export]
function array.get() {
    array.exists "$1" || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE
    local _array

    declare -n _array="$1"

    [ -n "${_array[$2]}" ] || return $COMMAND_FAILURE

    if [ -n "$3" ]; then
        declare -g "$3"="${_array[$2]}"
    else
        echo "${_array[$2]}"
    fi
    return $COMMAND_SUCCESS

}

# array.set <name> <offset> <value>
function array.set() {
    array.exists "$1" || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE
    [ -n "$3" ] || return $COMMAND_FAILURE
    local _array
    declare -n _array="$1"
    _array[$2]="$3"
    return $COMMAND_SUCCESS
}

# Map array and replaces value from function output
# array.map <name> <new_name> <function_name>
function array.map() {

    array.exists "$1" || return $COMMAND_FAILURE
    [ "$1" == "$2" ] && return $COMMAND_FAILURE
    declare -F "$3" &>>/dev/null || return $COMMAND_NOT_FOUND
    array.create "$2" || return $COMMAND_FAILURE

    local _array _new _key _value _newvalue

    declare -n _array="$1"
    declare -n _new="$2"

    for _key in "${!_array[@]}"; do
        _value="${_array[$_key]}"
        _newvalue="$($3 "$_value" "$_key")"
        if [ -n "$_newvalue" ]; then
            _new[$_key]="$_newvalue"
        else
            _new[$_key]="$_value"
        fi
    done

}
