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
