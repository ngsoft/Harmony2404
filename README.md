# Harmony with FLIRC Remote control Application 

That application includes a systemd daemon that write all the inputs from a [FLIRC USB Receiver](https://flirc.tv/products/flirc-usb-receiver) to a pipe file located `/var/run/flirc/flirc.pipe`

That application will include a [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API) that will read the pipe and send message events

A gui or console application will be created to read from that websocket and manage applications with the remote control using activities like [Remote Buddy](https://www.iospirit.com/products/remotebuddy/)


## Installation (Ubuntu LTS 24.04)

First install the flirc binaries (flirc_util, irtools) from https://flirc.com/ubuntu-software-installation-guide copying them into `/usr/local/bin` except for the appimage, the daemon uses `flirc_util` to check if the device is connected and to kick the device from bootloader when its power is interrupted.

Some extra packages are required to make it work (some packages have been renamed but its detected by apt)
```shell
sudo apt install libhidapi-hidraw0 libqt5core5a libqt5network5 libqt5xml5  libqt5xmlpatterns5 libhidapi-dev qtbase5-dev curl git
```

First, clone this repository

```shell
cd /tmp
git clone git@github.com:ngsoft/Harmony2404.git harmony
sudo mv harmony /opt
sudo chown -R $USER:$USER /opt/harmony
cd /opt/harmony
rm -rf .git
git init
git checkout -b main
git add .
git commit -m ":tada: first commit"
```

Some packages are required for the gui app to work, install them

```shell
sudo apt install xdotool evemu-tools
```

Some other packages can be installed if using gnome

```shell
sudo apt install libcanberra-gtk-module libcanberra-gtk3-module
```







Then install the service
The package is auto installable on first execution

```shell
sudo /opt/harmony/usr/local/bin/flircd
sudo systemctl daemon-reload
sudo systemctl start flircd
```

### InputLirc

We use a custom binary of [inputlirc](https://github.com/gsliepen/inputlirc) to capture the virtual keyboard
InputLirc source code is located at [lib/harmony/dist/inputlirc](./lib/harmony/dist/inputlirc)

### Logs

Logs are located in these files

```shell
tail -f /var/log/flircd.log
```

```shell
tail -f /var/log/flircws.log
```

### To read from the socket
```shell
sudo apt install netcat-traditional
```

```shell
nc -U /var/run/lirc/lircd
```

### The websocket

If flircd is running a websocket is running on port 9030

```url
ws://localhost.local:9030/ws
```

### The gui

For the gui app I will use [wails](https://wails.io/docs/introduction/) + [svelte](https://svelte.dev/docs/introduction)
