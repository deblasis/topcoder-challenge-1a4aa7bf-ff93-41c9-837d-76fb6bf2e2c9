package messaging

import (
	"errors"
	"sync"

	"github.com/deblasis/edgex-foundry-datamonitor/internal/config"
	edgexM "github.com/edgexfoundry/go-mod-messaging/v2/messaging"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
)

type Client struct {
	sync.Mutex
	edgeXClient edgexM.MessageClient
	cfg         *config.Config

	IsConnected  bool
	IsConnecting bool

	OnConnect func()
}

func NewClient(cfg *config.Config) (*Client, error) {

	cnf := config.DefaultConfig()

	if cfg != nil {
		cnf = cfg
	}

	c := &Client{
		Mutex:        sync.Mutex{},
		cfg:          cnf,
		IsConnected:  false,
		IsConnecting: false,
		OnConnect:    func() {},
	}

	return c, nil
}

func (c *Client) Connect() error {
	c.Lock()
	defer c.Unlock()

	c.IsConnecting = true
	defer func() {
		c.IsConnecting = false
	}()

	messageBus, err := edgexM.NewMessageClient(types.MessageBusConfig{
		SubscribeHost: types.HostInfo{
			Host:     config.StringVal(c.cfg.RedisHost),
			Port:     config.IntVal(c.cfg.RedisPort),
			Protocol: edgexM.Redis,
		},
		Type: edgexM.Redis,
	})

	if err != nil {
		return err
		//TODO log
		//LoggingClient.Error("failed to create messaging client: " + err.Error())
	}

	c.edgeXClient = messageBus

	err = c.edgeXClient.Connect()
	if err != nil {
		return err
	}

	c.IsConnected = true
	c.OnConnect()
	return nil
}

func (c *Client) Disconnect() error {
	c.Lock()
	defer c.Unlock()
	c.edgeXClient.Disconnect()
	c.IsConnected = false
	//TODO handle error
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
