class_loading=""

# Creates a new pseudo class instance
# new <constructor> <instance> ...[arguments]
function new() {
    [ -n "$1" ] || throw new: Invalid constructor
    [ -n "$2" ] || throw new: Invalid instance

    local self __CLASS__ __FILE__ ext _lib retval loaded classins baseins func push prev_loading
    __CLASS__="$1"
    self="$2"
    shift 2

    if [ "$__CLASS__" == "${class_loading}" ]; then
        throw "new: Cannot load $__CLASS__ inside $__CLASS__"
    fi

    declare -n loaded="${self}_instance"

    if [ -n "$loaded" ]; then
        ${self}.destroy
    fi

    if [ "$__CLASS__" == "Class" ]; then
        throw "new: Cannot instanciate abstract Class $__CLASS__"
    fi

    if ! declare -F "$__CLASS__" &>>/dev/null; then
        for ext in $(splitpath "${CLASS_EXT}"); do
            if __FILE__="$(script.find "${__CLASS__}.${ext}")"; then
                break
            fi
        done
        [ -n "$__FILE__" ] || throw "new: Cannot find class $__CLASS__"
        source "$__FILE__"
    fi

    if ! declare -F "$__CLASS__" &>>/dev/null; then
        throw "new: Cannot find constructor $__CLASS__"
    fi
    prev_loading="$class_loading"
    class_loading="$__CLASS__"

    # creates customized functions
    baseins="Class_${self}"
    classins="${__CLASS__}_${self}"

    # load Class
    func="$(declare -f Class)"
    source <(printf '%s' "${baseins}(${func#*\(}" | sed "s/self/$self/g")

    # load __CLASS__
    func="$(declare -f "$__CLASS__")"
    source <(printf '%s' "${classins}(${func#*\(}" | sed "s/self/$self/g")

    # execute constructor
    "$baseins"
    unset -f "$baseins"
    if [ -n "$__FILE__" ]; then
        push=true
        pushd "$(dirname "$__FILE__")" &>>/dev/null
    fi
    "$classins" "$@"
    unset -f "$classins"
    ${self}.constructor "$@"
    retval=$?
    unset -f "${self}.constructor"
    class_loading="$prev_loading"
    if [ "$push" == "true" ]; then
        popd &>>/dev/null
    fi
    return $retval
}

# Destroy an instance
# destroy <instance>
function destroy() {
    [ -n "$1" ] || return $COMMAND_FAILURE

    if declare -f "${1}.destroy" &>>/dev/null; then
        ${1}.destroy
        return $COMMAND_SUCCESS
    fi
    return $COMMAND_NOT_FOUND
}

# instanceof <instance> <constructor>
function instanceof() {
    [ -n "$1" ] || return $COMMAND_FAILURE
    [ -n "$2" ] || return $COMMAND_FAILURE
    declare -n _constructor="${1}_instance"
    [ "$_constructor" == "$2" ]
}

##
# Base Bash Class
# defines public getters/setters
# defines constructor/destructors
# every functions can be overriden inside the class
##
function Class {

    #class public properties
    declare -g -A self=()

    #used by instanceof
    self_instance="$__CLASS__"
    self_file="$__FILE__"

    ##
    # Holds reference to the class constructor name
    ##
    self.__CLASS__() {
        echo "$self_instance"
    }

    ##
    # Holds reference to the class filename
    ##
    self.__FILE__() {
        echo "$self_file"
    }

    ##
    # Create a property
    # usage: self.prop name = value || self.prop name value
    ##
    self.prop() {
        [ -n "$1" ] || throw self.prop: Invalid argument count
        if [ -z "$2" ]; then
            self.get "$1"
            return $?
        fi

        local _prop="$1"
        shift
        if [ "$1" == "=" ]; then shift; fi
        self.set "$_prop" "$1"
    }

    ##
    # Property getter
    # usage: self.get <name> [output]
    ##
    self.get() {
        [ -n "$1" ] || throw self.get: invalid argument count.
        [ -z "${self[$1]}" ] && return $COMMAND_FAILURE
        if [ -n "$2" ]; then
            declare -g "$2"="${self[$1]}"
        else
            echo "${self[$1]}"
        fi
        return $COMMAND_SUCCESS
    }

    ##
    # Property setter
    # usage: self.set <name> <value>
    ##
    self.set() {
        [ -n "$1" ] || throw self.set: invalid argument count.
        self[$1]="$2"
        return $COMMAND_SUCCESS
    }

    ##
    # Default class destructor
    # removes references to instance name
    ##
    self.destroy() {
        local method var
        unset -v self
        for method in $(compgen -A function self); do
            unset -f "${method}"
        done

        for var in $(compgen -A variable self); do
            unset -v "${var}"
        done
    }

    ##
    # Debug your script
    ##
    self.__dump() {
        local _res _params _method
        _params="$(declare -p self)"
        echo "class ${self_instance} (#self) {"

        echo "  properties: ${_params#*=}"

        for _res in $(compgen -A function self); do
            _method="$(declare -F "$_res")"
            echo "  ${_method}()"
        done
        # for _res in $(compgen -A variable self); do
        #     declare -p "$_res"
        # done

        echo }
    }

    ##
    # Default constructor
    ##
    self.constructor() {
        pass
    }
}
