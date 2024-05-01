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

Then install the Lirc suite and disable their daemons (we use a custom binary of [inputlirc](https://github.com/gsliepen/inputlirc) to capture the virtual keyboard) and irexec to read from the socket and write to the pipe

InputLirc source code is located at [lib/harmony/dist/inputlirc](./lib/harmony/dist/inputlirc)


```shell
sudo apt install lirc lirc-x xdotool libcanberra-gtk-module libcanberra-gtk3-module evemu-tools
for act in stop disable mask; do sudo systemctl $act lircd-setup; done
for act in stop disable mask; do sudo systemctl $act lircd-uinput; done
for act in stop disable mask; do sudo systemctl $act lircdm; done
for act in stop disable mask; do sudo systemctl $act irexec; done
for act in stop disable mask; do sudo systemctl $act lircd; done
for act in stop disable mask; do sudo systemctl $act lircd.socket; done
```

Then install the service

```shell
sudo cp /opt/harmony/etc/systemd/system/flirc.service /etc/systemd/system
sudo systemctl daemon-reload
sudo systemctl enable flirc.service 
sudo systemctl start flirc.service 
sudo systemctl status flirc.service 
```

To read from the pipe from the command line, run this command, then press a button on your remote control
```shell
if read input<"/var/run/flirc/flirc.pipe"; then echo $input; fi
```

### To read from the socket
```shell
sudo apt install netcat-traditional
```

```shell
nc -U /var/run/lirc/lircd
```