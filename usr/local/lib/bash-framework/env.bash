test -n "$_ROOT" || exit 1

##
# Main variables
##

if test -n "$INCLUDE_PATH"; then
    INCLUDE_PATH="${_ROOT}:${INCLUDE_PATH}"
else
    INCLUDE_PATH="${_ROOT}"
fi
REQUIRE_EXT=bash:sh
REQUIRE_LOADED=()
CLASS_EXT=class.bash
test -n "$LOGFILE" || LOGFILE=""
