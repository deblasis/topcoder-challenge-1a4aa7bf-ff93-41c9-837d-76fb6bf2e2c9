// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package services

import (
	"sync"

	"fyne.io/fyne/v2"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
)

type AppManager struct {
	sync.RWMutex
	client           *messaging.Client
	config           *config.Config
	currentContainer *fyne.Container

	navBar *fyne.Container

	ep *EventProcessor

	drawFn func(*fyne.Container)
}

func NewAppManager(client *messaging.Client, cfg *config.Config, ep *EventProcessor) *AppManager {
	return &AppManager{
		RWMutex: sync.RWMutex{},
		client:  client,
		config:  cfg,
		ep:      ep,
	}
}

func (a *AppManager) GetEventProcessor() *EventProcessor {
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
		if a.GetConnectionState() == ClientConnected {
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
		return ClientConnected
	}
	if a.client.IsConnecting {
		return ClientConnecting
	}
	return ClientDisconnected
}

func (a *AppManager) GetRedisHostPort() (string, int) {
	return a.config.GetRedisHost(), a.config.GetRedisPort()
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
