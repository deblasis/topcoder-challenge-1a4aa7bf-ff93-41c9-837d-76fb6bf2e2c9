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
	"testing"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/stretchr/testify/require"
)

func Test_IngestEvents(t *testing.T) {

	db := NewDB(1000)
	db.IngestEvent(dummyEvent())

	require.Equal(t, 1, db.readings.Count())
	require.Equal(t, 1, db.events.Count())

	// limiting to 2 records, adding 2 more
	db.UpdateBufferSize(2)

	db.IngestEvent(dummyEvent())
	db.IngestEvent(dummyEvent())

	// not filtering, we expect 2 events
	evts := db.GetEvents()
	require.Equal(t, 2, len(evts))

	// not filtering, we expect 2 readings
	rdngs := db.GetReadings()
	require.Equal(t, 2, len(rdngs))
	require.NotEmpty(t, rdngs[0].ResourceName)

	require.Equal(t, 2, db.readings.Count())
	require.Equal(t, 2, db.events.Count())

	db.IngestEvent(interestingEvent())
	require.Equal(t, 0, len(db.matchedEventIds.Serials))

	evts = db.GetEvents()
	require.Equal(t, 2, len(evts))

	db.UpdateFilter("interesting")
	require.Equal(t, 1, len(db.matchedEventIds.Serials))

	evts = db.GetEvents()
	require.Equal(t, 1, len(evts))
	require.Equal(t, "interesting_id", evts[0].Id)

	// only the reading should be matched
	db.UpdateFilter("interesting_reading_id")
	require.Equal(t, 0, len(db.matchedEventIds.Serials))

	// expecting not matches
	evts = db.GetEvents()
	require.Equal(t, 0, len(evts))

	require.Equal(t, 2, len(db.matchedReadingIds.Serials))

	// the readings should match a property on the parent event
	db.UpdateFilter("interesting_id")
	require.Equal(t, 2, len(db.matchedReadingIds.Serials))

	// no matches (it means that we return everything return all)
	db.UpdateFilter("")
	require.Equal(t, 0, len(db.matchedEventIds.Serials))
	require.Equal(t, 0, len(db.matchedReadingIds.Serials))

	// full buffer
	evts = db.GetEvents()
	require.Equal(t, 2, len(evts))

	// we keep only the last (interesting) event so we expect no matches
	db.UpdateBufferSize(1)

	// new buffer
	evts = db.GetEvents()
	require.Equal(t, 1, len(evts))

	db.UpdateFilter("event_id")
	require.Equal(t, 0, len(db.matchedEventIds.Serials))
	require.Equal(t, 0, len(db.matchedReadingIds.Serials))

	// no matches
	evts = db.GetEvents()
	require.Equal(t, 0, len(evts))

}

func dummyEvent() dtos.Event {
	return dtos.Event{
		Versionable: common.Versionable{},
		Id:          "event_id",
		DeviceName:  "device",
		ProfileName: "profile",
		Created:     1,
		Origin:      2,
		Readings: []dtos.BaseReading{
			{
				Versionable:  common.Versionable{},
				Id:           "reading_id",
				Created:      1,
				Origin:       2,
				DeviceName:   "device",
				ResourceName: "resource",
				ProfileName:  "profile",
				ValueType:    "Int32",
				BinaryReading: dtos.BinaryReading{
					BinaryValue: nil,
					MediaType:   "",
				},
				SimpleReading: dtos.SimpleReading{
					Value: "1151651",
				},
			},
		},
		Tags: map[string]string{},
	}
}

func interestingEvent() dtos.Event {
	return dtos.Event{
		Versionable: common.Versionable{},
		Id:          "interesting_id",
		DeviceName:  "interesting_device",
		ProfileName: "interesting_profile",
		Created:     1,
		Origin:      2,
		Readings: []dtos.BaseReading{
			{
				Versionable:  common.Versionable{},
				Id:           "interesting_reading_id",
				Created:      1,
				Origin:       2,
				DeviceName:   "interesting_device",
				ResourceName: "interesting_resource",
				ProfileName:  "interesting_profile",
				ValueType:    "Int32",
				BinaryReading: dtos.BinaryReading{
					BinaryValue: nil,
					MediaType:   "",
				},
				SimpleReading: dtos.SimpleReading{
					Value: "interesting_1151651",
				},
			},
			{
				Versionable:  common.Versionable{},
				Id:           "interesting_reading_id2",
				Created:      1,
				Origin:       2,
				DeviceName:   "interesting_device",
				ResourceName: "interesting_resource",
				ProfileName:  "interesting_profile",
				ValueType:    "Int32",
				BinaryReading: dtos.BinaryReading{
					BinaryValue: []byte{1, 2, 3, 4},
					MediaType:   "somebinary",
				},
				SimpleReading: dtos.SimpleReading{
					Value: "interesting_1151651",
				},
			},
		},
		Tags: map[string]string{},
	}
}
