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
package messaging

import (
	"errors"
	"log"
	"sync"

	"github.com/deblasis/edgex-foundry-datamonitor/config"
	edgexM "github.com/edgexfoundry/go-mod-messaging/v2/messaging"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
)

type Client struct {
	sync.Mutex
	edgeXClient edgexM.MessageClient
	cfg         *config.Config

	IsConnected  bool
	IsConnecting bool

	OnConnect func() bool
}

func NewClient(cfg *config.Config) (*Client, error) {

	c := &Client{
		Mutex:        sync.Mutex{},
		cfg:          cfg,
		IsConnected:  false,
		IsConnecting: false,
		OnConnect:    func() bool { return true },
	}

	return c, nil
}

func (c *Client) Connect() error {
	c.Lock()
	defer c.Unlock()
	c.IsConnected = false

	log.Printf("connecting to %v:%v\n", c.cfg.GetRedisHost(), c.cfg.GetRedisPort())

	c.IsConnecting = true
	defer func() {
		c.IsConnecting = false
	}()

	messageBus, err := edgexM.NewMessageClient(types.MessageBusConfig{
		SubscribeHost: types.HostInfo{
			Host:     c.cfg.GetRedisHost(),
			Port:     c.cfg.GetRedisPort(),
			Protocol: edgexM.Redis,
		},
		Type: edgexM.Redis,
	})

	if err != nil {
		log.Println(err)
		return err
	}

	c.edgeXClient = messageBus

	err = c.edgeXClient.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	// edgex doesn't return error on connect... but only on Suscribe / Publish
	// that's why we have to do something like this, which is not ideal
	c.IsConnected = c.OnConnect()

	return nil
}

func (c *Client) Disconnect() error {
	c.Lock()
	defer c.Unlock()
	c.edgeXClient.Disconnect()
	c.IsConnected = false
	return nil
}

func (c *Client) Subscribe(topic string) (chan types.MessageEnvelope, chan error) {

	errorChannel := make(chan error)
	if c.edgeXClient == nil {
		errorChannel <- errors.New("client not initialized")
		return nil, errorChannel
	}

	messages := make(chan types.MessageEnvelope)

	err := c.edgeXClient.Subscribe([]types.TopicChannel{
		{
			Topic:    topic,
			Messages: messages,
		},
	}, errorChannel)

	if err != nil {
		errorChannel <- err
	}

	return messages, errorChannel
}
