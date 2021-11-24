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
	"encoding/json"
	"runtime"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/asecurityteam/rolling"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type EventProcessor struct {
	eventsChannel <-chan *dtos.Event

	state chan processorState

	eventReceivedChannel   chan struct{}
	readingReceivedChannel chan struct{}

	TotalNumberEvents   int
	TotalNumberReadings int

	EventsPerSecondLastMinute   float64
	ReadingsPerSecondLastMinute float64

	LastEvents shortMemoryEventsSlicer

	eventListeners []EventListener

	sync.RWMutex
}

func NewEventProcessor(eventsChannel chan *dtos.Event) *EventProcessor {
	return &EventProcessor{
		eventsChannel: eventsChannel,

		state: make(chan processorState, 1),

		eventListeners: make([]EventListener, 0),

		eventReceivedChannel:   make(chan struct{}, config.MaxBufferSize),
		readingReceivedChannel: make(chan struct{}, config.MaxBufferSize),
		LastEvents:             newTopNEventSlicer(5),
	}
}

func (ep *EventProcessor) Activate() {
	ep.Lock()
	defer ep.Unlock()
	ep.state <- Running
}

func (ep *EventProcessor) Deactivate() {
	ep.Lock()
	defer ep.Unlock()
	ep.state <- Paused
}

func (ep *EventProcessor) AttachListener(listener EventListener) {
	ep.eventListeners = append(ep.eventListeners, listener)
}

func (ep *EventProcessor) processEvent(event *dtos.Event) {

	for _, listener := range ep.eventListeners {
		listener.OnEventReceived(*event)
	}

	ep.eventReceivedChannel <- struct{}{}
	ep.TotalNumberEvents++
	ep.TotalNumberReadings += len(event.Readings)

	ep.LastEvents.Add(event)

	for range event.Readings {
		ep.readingReceivedChannel <- struct{}{}
	}
}

func (ep *EventProcessor) Run() {

	timeWindow := rolling.NewWindow(1000 * 60)
	rollingEventsCounter := rolling.NewTimePolicy(timeWindow, time.Millisecond)
	rollingReadingsCounter := rolling.NewTimePolicy(timeWindow, time.Millisecond)

	go func() {
		for range ep.eventReceivedChannel {
			// if we get more events than we can process, the Append below panics, catching it, needs investigation
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic occurred:", err)
				}
			}()
			rollingEventsCounter.Append(1)
		}
	}()

	go func() {
		for range ep.readingReceivedChannel {
			// if we get more events than we can process, the Append below panics, catching it, needs investigation
			defer func() {
				if err := recover(); err != nil {
					log.Errorf("panic occurred:", err)
				}
			}()
			rollingReadingsCounter.Append(1)
		}
	}()

	go func() {
		for range time.Tick(time.Millisecond * 200) {
			eventsPerMinute := rollingEventsCounter.Reduce(rolling.Sum)
			ep.EventsPerSecondLastMinute = eventsPerMinute / 60
		}
	}()

	go func() {
		for range time.Tick(time.Millisecond * 200) {
			readingsPerMinute := rollingReadingsCounter.Reduce(rolling.Sum)
			ep.ReadingsPerSecondLastMinute = readingsPerMinute / 60
		}
	}()

	state := Running

	for {

		select {

		case state = <-ep.state:
			switch state {
			case Stopped:
				log.Info("EventsProcessor: Stopped")
				return
			case Running:
				log.Info("EventsProcessor: Running")
			case Paused:
				log.Info("EventsProcessor: Paused")
			}

		case event := <-ep.eventsChannel:

			runtime.Gosched()

			if state == Paused {
				break
			}

			ep.processEvent(event)
		}
	}
}

type topNEvents struct {
	n      int
	events []*dtos.Event
}

func newTopNEventSlicer(n int) *topNEvents {
	return &topNEvents{
		n:      n,
		events: []*dtos.Event{},
	}
}

func (t *topNEvents) Add(e *dtos.Event) {
	if len(t.events) == t.n {
		t.events = t.events[1:]
	}

	t.events = append(t.events, e)
}

func (l *topNEvents) Get() []*dtos.Event {
	return l.events
}

func (l *topNEvents) GetJson() string {
	j, _ := json.Marshal(l.events)
	return string(j)
}

type shortMemoryEventsSlicer interface {
	Add(e *dtos.Event)
	Get() []*dtos.Event
	GetJson() string
}

type EventListeners []EventListener
type EventListener interface {
	OnEventReceived(event dtos.Event)
}
