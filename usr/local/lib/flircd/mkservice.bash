unit="/etc/systemd/system/${APP_NAME}.service"

if [ ! -e "$unit" ]; then

    echo "Installing service ..." >&2

    cat >"$unit" <<EOL
[Unit]
Description=FLIRC Manager
After=network.target

[Service]
Type=simple
Environment="FLIRC_SYSTEMD=1"
ExecStart=${FLIRC_PREFIX}/usr/local/bin/flircd
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
EOL
    systemctl enable "${APP_NAME}.service"
    echo "Service ${APP_NAME}.service has been enabled, please run:" >&2
    echo "systemctl start "${APP_NAME}.service"" >&2
    exit $COMMAND_SUCCESS
fi
