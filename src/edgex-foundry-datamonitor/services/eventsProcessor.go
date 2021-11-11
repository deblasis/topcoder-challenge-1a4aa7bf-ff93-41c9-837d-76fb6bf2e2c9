package services

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/asecurityteam/rolling"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type EventProcessor struct {
	eventsChannel <-chan *dtos.Event

	state chan processorState

	lastEventChannel   chan dtos.Event
	lastReadingChannel chan dtos.BaseReading

	TotalNumberEvents   int
	TotalNumberReadings int

	EventsPerSecondLastMinute   float64
	ReadingsPerSecondLastMinute float64

	LastEvents shortMemoryEventsSlicer

	sync.RWMutex
}

func NewEventProcessor(eventsChannel chan *dtos.Event) *EventProcessor {
	return &EventProcessor{
		eventsChannel: eventsChannel,

		state: make(chan processorState, 1),

		lastEventChannel:   make(chan dtos.Event, 1),
		lastReadingChannel: make(chan dtos.BaseReading, 1),
		LastEvents:         newTopNEventSlicer(5),
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

func (ep *EventProcessor) processEvent(event *dtos.Event) {

	ep.lastEventChannel <- *event
	ep.TotalNumberEvents++
	ep.TotalNumberReadings += len(event.Readings)

	ep.LastEvents.Add(event)

	for _, reading := range event.Readings {
		ep.lastReadingChannel <- reading
	}
}

func (ep *EventProcessor) Run() {

	timeWindow := rolling.NewWindow(1000 * 60)

	var rollingEventsCounter = rolling.NewTimePolicy(timeWindow, time.Millisecond)
	go func() {
		for range time.Tick(time.Millisecond) {
			<-ep.lastEventChannel
			rollingEventsCounter.Append(1)
		}
	}()

	var rollingReadingsCounter = rolling.NewTimePolicy(timeWindow, time.Millisecond)
	go func() {
		for range time.Tick(time.Millisecond) {
			<-ep.lastReadingChannel
			rollingReadingsCounter.Append(1)
		}
	}()

	go func() {
		for {
			eventsPerMinute := rollingEventsCounter.Reduce(rolling.Sum)
			ep.EventsPerSecondLastMinute = eventsPerMinute / 60
			time.Sleep(time.Millisecond * 200) //throttling a bit
		}
	}()

	go func() {
		for {
			readingsPerMinute := rollingReadingsCounter.Reduce(rolling.Sum)
			ep.ReadingsPerSecondLastMinute = readingsPerMinute / 60
			time.Sleep(time.Millisecond * 200) //throttling a bit
		}
	}()

	state := Running

	for {

		select {

		case state = <-ep.state:
			switch state {
			case Stopped:
				log.Println("EventsProcessor: Stopped")
				return
			case Running:
				log.Println("EventsProcessor: Running")
			case Paused:
				log.Println("EventsProcessor: Paused")
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
