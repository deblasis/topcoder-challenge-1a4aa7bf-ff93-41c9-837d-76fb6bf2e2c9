package messaging

import (
	"encoding/json"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type eventEnvelope struct {
	Event *dtos.Event `json:"event"`
}

func ParseEvent(eventPayloadBytes []byte) (*dtos.Event, error) {
	e := &eventEnvelope{}
	if err := json.Unmarshal(eventPayloadBytes, e); err != nil {
		return nil, err
	}
	return e.Event, nil
}
