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

package messaging_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
	"github.com/deblasis/edgex-foundry-datamonitor/mocks"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-messaging/v2/pkg/types"
	"github.com/stretchr/testify/require"
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

			mockApp := new(mocks.MockApp)
			mockPreferences := new(mocks.MockPreferences)

			mockPreferences.On("StringWithFallback", "_RedisHost", "localhost").Return("localhost")
			mockPreferences.On("IntWithFallback", "_RedisPort", 6379).Return(6379)

			mockApp.On("Preferences").Return(mockPreferences)

			c := newDefaultClient(mockApp)

			events := make(chan *dtos.Event)

			var (
				errs     = make(chan error)
				messages = make(chan types.MessageEnvelope)
			)

			c.OnConnect = func() bool {
				messages, errs = c.Subscribe(config.DefaultEventsTopic)

				return true
			}

			err := c.Connect()
			require.Nil(t, err)

			ep := services.NewEventProcessor(events)
			go ep.Run()

			gracePeriod := time.NewTimer(10 * time.Second)
		LOOP:
			for {
				select {

				case e := <-errs:
					t.Fatalf("Client.Subscribe() got error = %v", e)

				case msgEnvelope := <-messages:
					gracePeriod.Stop()
					event, err := messaging.ParseEvent(msgEnvelope.Payload)
					events <- event

					if err != nil {
						t.Fatalf("Client.Subscribe() got error while parsing Event = %v", err)
					}

					t.Log(event)
					t.Logf("Client.Subscribe() got msgEnvelope = %v", msgEnvelope)
					break LOOP
				case <-gracePeriod.C:
					ep.Deactivate()
					t.Fatal("Client.Subscribe() no messages received within the gracePeriod")
				}
			}
		})
	}
}

func newDefaultClient(app fyne.App) *messaging.Client {

	cfg := config.GetConfig(app)

	c, _ := messaging.NewClient(cfg)
	return c
}
