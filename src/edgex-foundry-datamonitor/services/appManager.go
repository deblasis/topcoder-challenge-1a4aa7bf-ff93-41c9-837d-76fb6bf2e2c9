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
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
)

type AppManager struct {
	sync.RWMutex
	client           *messaging.Client
	config           *config.Config
	currentContainer *fyne.Container

	CurrentPage widget.TreeNodeID

	navBar *fyne.Container

	db *DB
	ep *EventProcessor

	pageHandlers map[widget.TreeNodeID]PageHandler

	drawFn func(*fyne.Container)

	sessionState *SessionState
}

func NewAppManager(client *messaging.Client, cfg *config.Config, ep *EventProcessor, db *DB) *AppManager {

	return &AppManager{
		RWMutex: sync.RWMutex{},
		client:  client,
		config:  cfg,
		db:      db,
		ep:      ep,

		pageHandlers: make(map[widget.TreeNodeID]PageHandler),

		sessionState: &SessionState{},
	}
}

func (a *AppManager) SetPageHandler(page widget.TreeNodeID, handler PageHandler) {
	a.pageHandlers[page] = handler
}

func (a *AppManager) GetPageHandler(page widget.TreeNodeID) PageHandler {
	return a.pageHandlers[page]
}

func (a *AppManager) GetEventProcessor() *EventProcessor {
	return a.ep
}

func (a *AppManager) GetDB() *DB {
	return a.db
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

type SessionState struct {
	DataPage_SelectedDataType *string
	DataPage_Search           *string
	DataPage_BufferSize       *int
}

func (a *AppManager) SetDataPageSelectedDataType(dt string) {
	a.Lock()
	defer a.Unlock()
	a.sessionState.DataPage_SelectedDataType = config.String(dt)
}

func (a *AppManager) SetDataPageSearch(search string) {
	a.Lock()
	defer a.Unlock()
	a.sessionState.DataPage_Search = config.String(search)
	a.db.UpdateFilter(search)
}

func (a *AppManager) SetDataPageBufferSize(bs int) {
	a.Lock()
	defer a.Unlock()
	a.sessionState.DataPage_BufferSize = config.Int(bs)
	a.db.UpdateBufferSize(int64(bs))
}

func (a *AppManager) GetDataPageSelectedDataType() *string {
	a.RLock()
	defer a.RUnlock()
	return a.sessionState.DataPage_SelectedDataType
}

func (a *AppManager) GetDataPageSearch() *string {
	a.RLock()
	defer a.RUnlock()
	return a.sessionState.DataPage_Search
}

func (a *AppManager) GetDataPageBufferSize() *int {
	a.RLock()
	defer a.RUnlock()
	return a.sessionState.DataPage_BufferSize
}
