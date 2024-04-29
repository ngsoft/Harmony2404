#
# Process Managemant class
#
function Process(){

    this.constructor(){
        ! is.empty $1 && this.process $1
    }

    this.process(){
        if [ "$1" == "=" ] && ! is.empty $2; then
            this.prop process "$2"
        elif ! is.empty $1; then this.prop process "$1"
        else this.prop process; fi
    }

    this.wait.running(){
        is.empty ${_this[process]} && return 1
        local timeout=10 cnt=0
        is.int.unsigned $1 && timeout=$1
        # loops are 200ms
        ((timeout=timeout*5))
        while true; do
            sleep .2
            ((cnt+=1))
            this.running && return 0
            [[ $cnt -gt $timeout ]] && return 1
        done
    }

    this.ispid(){
        is.int.unsigned ${_this[process]}
    }

    this.running(){
        is.empty ${_this[process]} && return 1
        if this.ispid; then cmd ps "${_this[process]}"
        else cmd pgrep "${_this[process]}"; fi
    }

    this.kill(){
        is.empty ${_this[process]} && return 1
        if this.ispid; then cmd kill ${_this[process]}
        else cmd killall "${_this[process]}"; fi
    }

    this.getpid(){
        is.empty ${_this[process]} && return 1
        if this.ispid; then echo ${_this[process]}
        else pgrep -n "${_this[process]}"; fi
    }

    this.getpids(){
        is.empty ${_this[process]} && return 1
        if this.ispid; then echo ${_this[process]}
        else pgrep "${_this[process]}"; fi
    }

    this.countpids(){
        is.empty ${_this[process]} && return 1
        if this.ispid; then
            if this.running; then echo 1
            else echo 0; fi
        else pgrep -c "${_this[process]}"; fi
    }

    this.signal(){
        is.empty $1 && return 1
        local pid
        for pid in $(this.getpids); do
            cmd ps $pid && kill -s $1 $pid
        done
    }

}
