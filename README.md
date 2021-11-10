# topcoder-challenge-4c8811c4-6504-40a1-a9d0-ad25ee7c1af7
EdgeX Foundry(TM) - Build Simple Data Monitor UI #01



# Linux (ubuntu)
## https://developer.fyne.io/started/
sudo apt-get install golang gcc libgl1-mesa-dev xorg-dev


go install github.com/fyne-io/fyne-cross@latest

fyne-cross linux


https://github.com/edgexfoundry/edgex-go#zeromq
brew install zeromq


## IMPORTANT ZeroMq deprecation!
https://github.com/edgexfoundry/go-mod-messaging/issues/73

## WSL

export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2; exit;}'):0.0
export LIBGL_ALWAYS_INDIRECT=0


https://i.stack.imgur.com/4n4XH.png

Then, start a new instance of VcxSrv with and unselect the Native opengl box on the Extra Settings page, and select the Disable access control box. VcxSrv Extra Settings
