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
	edgeXClient  edgexM.MessageClient
	IsConnected  bool
	IsConnecting bool
}

func NewClient(cfg *config.Config) (*Client, error) {
	messageBus, err := edgexM.NewMessageClient(types.MessageBusConfig{
		PublishHost: types.HostInfo{
			Host:     config.StringVal(cfg.RedisHost),
			Port:     config.IntVal(cfg.RedisPort),
			Protocol: edgexM.Redis,
		},
		Type: edgexM.Redis})

	if err != nil {
		return nil, err
		//TODO log
		//LoggingClient.Error("failed to create messaging client: " + err.Error())
	}

	c := &Client{
		Mutex:        sync.Mutex{},
		edgeXClient:  messageBus,
		IsConnected:  false,
		IsConnecting: false,
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

	if c.edgeXClient == nil {
		return errors.New("client not initialized") //TODO refactor
	}

	err := c.edgeXClient.Connect()
	if err != nil {
		return err
	}

	c.IsConnected = true
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
