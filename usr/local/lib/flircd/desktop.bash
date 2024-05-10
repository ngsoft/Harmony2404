function desktop.running() {
    local c checks="plasmashell gnome-shell lxqt-session xfce4-session gnome-session-binary"

    for c in $checks; do
        if process.running "$c"; then
            return $COMMAND_SUCCESS
        fi
    done
    return $COMMAND_FAILURE
}

function sudo.install() {
    if [ ! -e "/etc/sudoers.d/flirc" ]; then
        echo "root ALL=(ALL:ALL) NOPASSWD:ALL" >/etc/sudoers.d/flirc
    fi
}

function session.detect() {
    if ! desktop.running; then return $COMMAND_FAILURE; fi
    if [ -z "$DISPLAY" ]; then
        log "INFO: display server is not available for current session"
        local info user disp
        for info in $(who); do
            if [ -z "$user" ]; then
                user="$info"
                continue
            fi
            if [ "$info" != "${info#(}" ]; then
                info="${info#(}"
                info="${info%)}"
                if [ -z "${info#:[0-9]}" ]; then
                    disp="$info"
                    break
                else
                    user=""
                fi
            fi
        done
        if [ -z "$user" ] || [ -z "$disp" ]; then return $COMMAND_FAILURE; fi
        export SESSION_USER="$user"
        export DISPLAY="$disp"
        log "INFO: found session user ${SESSION_USER}${DISPLAY}"
    fi
    return $COMMAND_SUCCESS
}

function sudo.execute() {
    [ -z "$1" ] && return $COMMAND_FAILURE
    if [ -n "$SESSION_USER" ]; then
        sudo.install
        sudo -u $SESSION_USER $@
    else
        $@
    fi
}
