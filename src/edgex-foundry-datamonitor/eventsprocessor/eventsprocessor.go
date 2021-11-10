package eventsprocessor

import (
	"fmt"
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

	// readingsChannel chan models.Reading

	TotalNumberEvents   int
	TotalNumberReadings int

	EventsPerSecondLastMinute   float64
	ReadingsPerSecondLastMinute float64

	LastDeviceNames      lastFiveStrings
	LastOriginTimestamps lastFiveInt

	sync.RWMutex
}

func New(eventsChannel chan *dtos.Event) *EventProcessor {
	return &EventProcessor{
		eventsChannel: eventsChannel,

		state: make(chan processorState, 1),

		lastEventChannel:   make(chan dtos.Event, 1),
		lastReadingChannel: make(chan dtos.BaseReading, 1),

		LastDeviceNames: lastFiveStrings{
			arr: make([]string, 0),
		},
		LastOriginTimestamps: lastFiveInt{
			arr: make([]int64, 0),
		},
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

// func (ep *EventProcessor) IsActive() bool {
// 	ep.RLock()
// 	defer ep.RUnlock()
// 	return ep.isActive
// }

func (ep *EventProcessor) processEvent(event *dtos.Event) {

	ep.lastEventChannel <- *event
	ep.TotalNumberEvents++
	ep.TotalNumberReadings += len(event.Readings)

	ep.LastDeviceNames.Add(event.DeviceName)
	ep.LastOriginTimestamps.Add(event.Origin)

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
				fmt.Println("Stopped")
				return
			case Running:
				fmt.Println("Running")
			case Paused:
				fmt.Println("Paused")
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

type processorState int

const (
	Stopped processorState = iota
	Paused
	Running
)

//TODO refactor
type lastFiveStrings struct {
	arr []string
}

func (l *lastFiveStrings) Add(a string) {
	if len(l.arr) == 5 {
		l.arr = l.arr[1:]
	}

	l.arr = append(l.arr, a)
}

func (l *lastFiveStrings) Get() []string {
	return l.arr
}

type lastFiveInt struct {
	arr []int64
}

func (l *lastFiveInt) Add(a int64) {
	if len(l.arr) == 5 {
		l.arr = l.arr[1:]
	}

	l.arr = append(l.arr, a)
}

func (l *lastFiveInt) Get() []int64 {
	return l.arr
}
