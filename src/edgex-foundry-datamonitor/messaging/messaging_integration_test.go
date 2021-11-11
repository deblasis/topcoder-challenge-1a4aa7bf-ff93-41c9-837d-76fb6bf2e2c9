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
//g o : b u i l d   i n t e g r ation
// b u i ld   i  ntegration

package messaging

import (
	"testing"
	"time"

	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/eventsprocessor"
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
