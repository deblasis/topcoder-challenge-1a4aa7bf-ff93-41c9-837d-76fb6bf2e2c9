package state

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/eventsprocessor"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
)

type AppManager struct {
	sync.RWMutex
	client           *messaging.Client
	config           *config.Config
	currentContainer *fyne.Container

	navBar *fyne.Container

	ep *eventsprocessor.EventProcessor

	drawFn func(*fyne.Container)
}

func NewAppManager(client *messaging.Client, cfg *config.Config, ep *eventsprocessor.EventProcessor) *AppManager {
	return &AppManager{
		RWMutex: sync.RWMutex{},
		client:  client,
		config:  cfg,
		ep:      ep,
	}
}

func (a *AppManager) GetEventProcessor() *eventsprocessor.EventProcessor {
	return a.ep
}

func (a *AppManager) SetCurrentContainer(container *fyne.Container, drawFn func(*fyne.Container)) {
	a.Lock()
	defer a.Unlock()
	a.currentContainer = container
	a.drawFn = drawFn
}

func (a *AppManager) SetNav(nav *fyne.Container) {
	a.Lock()
	defer a.Unlock()
	a.navBar = nav
}

func (a *AppManager) Refresh() {
	refreshNavBar := func() {
		if a.navBar == nil {
			return
		}
		if a.GetConnectionState() == Connected {
			a.navBar.Objects[1].(*fyne.Container).Objects[0].Show()
		} else {
			a.navBar.Objects[1].(*fyne.Container).Objects[0].Hide()
		}
		a.navBar.Refresh()
	}
	refreshContent := func() {
		if a.drawFn != nil && a.currentContainer != nil {
			a.drawFn(a.currentContainer)
		}
	}

	refreshNavBar()
	refreshContent()

}

func (a *AppManager) GetCurrentContainer() (*fyne.Container, func(*fyne.Container)) {
	a.RLock()
	defer a.RUnlock()
	return a.currentContainer, a.drawFn
}

// func (a *AppState) IsConnected() bool {
// 	a.RLock()
// 	defer a.RUnlock()
// 	return a.client.IsConnected
// }

// func (a *AppState) IsConnecting() bool {
// 	a.RLock()
// 	defer a.RUnlock()
// 	return a.client.IsConnecting
// }

func (a *AppManager) GetConnectionState() ConnectionState {
	a.RLock()
	defer a.RUnlock()
	if a.client.IsConnected {
		return Connected
	}
	if a.client.IsConnecting {
		return Connecting
	}
	return Disconnected
}

func (a *AppManager) GetRedisHostPort() (string, int) {
	a.RLock()
	defer a.RUnlock()

	if a.config.RedisHost == nil && a.config.RedisPort == nil {
		return config.RedisDefaultHost, config.RedisDefaultPort
	}

	return config.StringVal(a.config.RedisHost), config.IntVal(a.config.RedisPort)
}

func (a *AppManager) Connect() error {
	a.Lock()
	defer a.Unlock()
	a.ep.Activate()
	return a.client.Connect()
}

func (a *AppManager) Disconnect() error {
	a.Lock()
	defer a.Unlock()
	a.ep.Deactivate()
	return a.client.Disconnect()
}

type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connecting
	Connected
)
