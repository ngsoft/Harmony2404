#
# Contains Environment variables
#

if [ -z "$HARMONY_APPLICATION_PREFIX" ]; then
    echo '$HARMONY_APPLICATION_PREFIX is not set, please setup.' >&2
    exit 1
fi

DISPLAY=:0
HARMONY_LIB=$HARMONY_APPLICATION_PREFIX/lib/harmony
HARMONY_DIST=$HARMONY_LIB/dist
HARMONY_CONF=$HARMONY_APPLICATION_PREFIX/etc
HARMONY_SKEL=$HARMONY_LIB/skel
HARMONY_INCLUDE_PATH=$HARMONY_APPLICATION_PREFIX/bin:$HARMONY_LIB:$HARMONY_LIB/classes

# FLIRC Daemon Config
HARMONY_FLIRC_PATH="/var/run/flirc"
HARMONY_FLIRC_FIFO=$HARMONY_FLIRC_PATH/flirc.pipe
HARMONY_FLIRC_KEYMAPS=$HARMONY_CONF/flirc.keymaps
HARMONY_FLIRCD_PIDFILE=$HARMONY_FLIRC_PATH/flircd.pid
HARMONY_FLIRCD_IREXEC=$HARMONY_FLIRC_PATH/irexec.cfg
HARMONY_FLIRC_CONFIG=$HARMONY_CONF/flirc.conf
HARMONY_FIRCD_LOGFILE=$HARMONY_FLIRC_PATH/flircd.log

# defining PATH for easy access to libraries for plugins and activities scripts
PATH=$PATH:$HARMONY_APPLICATION_PREFIX/bin:$HARMONY_LIB

# Load this file only once
HARMONY_ENV_LOADED=1
