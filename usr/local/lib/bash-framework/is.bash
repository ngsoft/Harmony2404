##
# Comparison functions
##
function is() {
    [ "$1" == "$2" ]
}

function is.bool() {
    [[ "$1" =~ ^(true|false|TRUE|FALSE)$ ]]
}
function is.true() {
    [ "$1" == "true" ] || [ "$1" == "TRUE" ]
}
function is.int() {
    [[ "$1" =~ ^[\+\-]?[0-9]+$ ]]
}
function is.unsigned() {
    [[ "$1" =~ ^[0-9]+ ]]
}
function is.float() {
    [[ "$1" =~ ^[\+\-]?[0-9]*\.[0-9]+$ ]]
}

function is.empty() {
    [ -z "$1" ]
}
function is.function() {
    declare -F "$1" &>>/dev/null
}

function is.array() {
    [[ "$(declare -p "$1" 2>/dev/null)" =~ declare[^-]-[aA] ]]
}

function is.instance() {
    declare -p "${1}_instance" &>>/dev/null
}
