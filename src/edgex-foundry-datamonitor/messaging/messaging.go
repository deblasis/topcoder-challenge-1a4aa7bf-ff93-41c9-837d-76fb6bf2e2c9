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
		errorChannel <- errors.New("client not initialized") //TODO refactor
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

	// events := make(chan *dtos.Event)
	// msgEnvelope := <-messages
	// //gracePeriod.Stop()
	// event, err := parseEvent(msgEnvelope.Payload)
	// events <- event

	return messages, errorChannel
}
