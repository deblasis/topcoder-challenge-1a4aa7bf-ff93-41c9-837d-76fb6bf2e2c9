module github.com/deblasis/edgex-foundry-datamonitor

go 1.16

require fyne.io/fyne/v2 v2.1.1

require (
	github.com/asecurityteam/rolling v2.0.4+incompatible
	github.com/edgexfoundry/go-mod-core-contracts v0.1.149
	github.com/edgexfoundry/go-mod-messaging/v2 v2.0.1
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-gl/gl v0.0.0-20211025173605-bda47ffaa784 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20211024062804-40e447a793be // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/kelindar/column v0.0.0-20211106170543-f720749ebf55
	github.com/sirupsen/logrus v1.8.1
	github.com/srwiley/oksvg v0.0.0-20211104221756-aeb4ca2c1505 // indirect
	github.com/srwiley/rasterx v0.0.0-20210519020934-456a8d69b780 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/yuin/goldmark v1.4.2 // indirect
	golang.org/x/image v0.0.0-20211028202545-6944b10bf410 // indirect
	golang.org/x/net v0.0.0-20211105192438-b53810dc28af // indirect
	golang.org/x/sys v0.0.0-20211107104306-e0b2ad06fe42 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/edgexfoundry/go-mod-messaging/v2 v2.0.1 => github.com/deblasis/go-mod-messaging/v2 v2.0.2-dev.7.0.20211109140946-fc6f30b82981
