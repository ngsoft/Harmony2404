function Daemon {
    self.constructor() {
        is.empty $1 && throw "Cannot instantiate Daemon, no process provided"
        self[process]="$1"
    }

    self.pid() {
        process.pid "${self[process]}"
    }

    self.pids() {
        pgrep "${self[process]}"
    }

    self.count() {
        process.count "${self[process]}"
    }

    self.signal() {
        is.empty $1 && return $COMMAND_FAILURE
        local pid
        for pid in $(self.pids); do
            @ ps $pid && kill -s $1 $pid
        done
    }

    self.kill() {
        process.kill "${self[process]}"
    }

    self.running() {
        process.running "${self[process]}"
    }

    self.wait() {
        process.wait "${self[process]}" $@
    }
}

function ProcessList {

    self.constructor() {
        is.empty $1 && throw "Cannot instantiate self, no process provided"
        local _process _index=0
        array.create.list self_list $@
        for _process in "${self_list[@]}"; do
            new Daemon "selfprocess${_index}" "${_process}"
            ((_index += 1))
        done
    }

    self.pid() {
        local _process _pid
        for _process in "${!self_list[@]}"; do
            _pid=$(selfprocess${_process}.pid)
            if [ -n "$_pid" ]; then
                echo "$_pid"
                return $COMMAND_SUCCESS
            fi
        done
        return $COMMAND_FAILURE
    }

    self.pids() {
        local _process _result=$COMMAND_FAILURE
        for _process in "${!self_list[@]}"; do
            if selfprocess${_process}.pids; then
                _result=$COMMAND_SUCCESS
            fi
        done
        return $_result
    }

    self.count() {
        local _cnt=0 _add=0
        for _process in "${!self_list[@]}"; do
            _add="$(selfprocess${_process}.count "${_process}")"
            ((_cnt += _add))
        done
        echo "${_cnt}"
        [ "$_cnt" != "0" ]
    }

    self.signal() {
        is.empty $1 && return $COMMAND_FAILURE
        for _process in "${!self_list[@]}"; do
            selfprocess${_process}.signal $1
        done
    }

    self.kill() {
        for _process in "${!self_list[@]}"; do
            selfprocess${_process}.kill
        done
    }

    self.running() {
        @ self.pid
    }

    self.destroy() {
        local method var _process

        # remove Daemon instances first
        for _process in "${!self_list[@]}"; do
            selfprocess${_process}.destroy
        done

        unset -v self
        for method in $(compgen -A function self.); do
            unset -f "${method}"
        done

        for var in $(compgen -A variable self_); do
            unset -v "${var}"
        done
    }
}
