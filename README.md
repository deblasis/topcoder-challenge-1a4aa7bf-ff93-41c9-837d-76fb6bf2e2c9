# topcoder-challenge-1a4aa7bf-ff93-41c9-837d-76fb6bf2e2c9
EdgeX Foundry(TM) - Build Simple Data Monitor UI #02

TLDR; left is Linux, right is Windows
<img src="./assets/crossplatform.jpg" alt="crossplatform" />


## Connected state
<img src="./assets/connected.jpg" alt="connected" />


## IMPORTANT ZeroMq deprecation!

I had to patch the referenced  https://github.com/edgexfoundry/go-mod-messaging library because it uses a library that made me lose a whole day while trying to make it work in my environment. It will be soon deprecated as stated here https://github.com/edgexfoundry/go-mod-messaging/issues/73

Also, we use Redis, not zeromq so it’s totally fine

Currently my patched version is referenced with a go mod replace entry

<img src="./assets/gomodpatch.jpg" alt="go mod patch" />


## IMPORTANT Darwin compilation

❗❗❗❗ 🍎
>OSX/Darwin/Apple cross-compiling requires a darwin host and/or some manual steps along with the acceptance of Xcode license terms
Please follow the link below:
https://github.com/fyne-io/fyne-cross#build-the-docker-image-for-osxdarwinapple-cross-compiling

❗❗❗❗


## WSL2 GUI

In order to run the Linux version on WSL you need an X11 server like **VcXsrv** launched with the following settings:

<img src="https://i.stack.imgur.com/4n4XH.png" />

And then run the following command, at least this is what I had to do on my machine 😉

```bash
export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2; exit;}'):0.0
export LIBGL_ALWAYS_INDIRECT=0
```


## docker-compose

The docker compose includes the UI which is accessible here:

http://localhost:4000/#/dashboard

you can try to increase/decrease the amount of events being produced by changing the settings in here: http://localhost:4000/#/metadata/device-center/device-list


Or going directly into one of the devices' settings:
http://localhost:4000/#/metadata/device-center/edit-device?deviceName=Random-Float-Device


### License

Apache License Version 2.0