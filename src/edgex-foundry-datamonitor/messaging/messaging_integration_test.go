//g o : b u i l d   i n t e g r ation
// b u i ld   i  ntegration

package messaging

import (
	"testing"
	"time"

	"github.com/deblasis/edgex-foundry-datamonitor/eventsprocessor"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/config"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
)

func TestClient_Subscribe(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name  string
		args  args
		want  chan types.MessageEnvelope
		want1 chan error
	}{
		{
			name: "subscribe to events",
			args: args{
				topic: config.DefaultEventsTopic,
			},
			want:  make(chan types.MessageEnvelope),
			want1: make(chan error),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newDefaultClient()

			c.Connect()
			messages, errors := c.Subscribe(tt.args.topic)

			events := make(chan *dtos.Event)

			ep := eventsprocessor.New(events)
			go ep.Run()

			gracePeriod := time.NewTimer(10 * time.Second)

			for {
				select {
				//return

				case e := <-errors:
					t.Fatalf("Client.Subscribe() got error = %v", e)

				case msgEnvelope := <-messages:
					//gracePeriod.Stop()
					event, err := ParseEvent(msgEnvelope.Payload)
					events <- event

					if err != nil {
						t.Fatalf("Client.Subscribe() got error while parsing Event = %v", err)
					}
					t.Log(event)
					t.Logf("Client.Subscribe() got msgEnvelope = %v", msgEnvelope)
				case <-gracePeriod.C:
					//t.Fatal("Client.Subscribe() couldn't get any message within the gracePeriod")
					ep.Deactivate()
				}
			}
		})
	}
}

func newDefaultClient() *Client {
	c, _ := NewClient(nil)
	return c
}
