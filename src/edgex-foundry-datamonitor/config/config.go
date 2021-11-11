package config

import "fyne.io/fyne/v2"

const (
	PrefRedisHost = "_RedisHost"
	PrefRedisPort = "_RedisPort"

	PrefShouldConnectAtStartup        = "_ShouldConnectAtStartup"
	PrefEventsTableSortOrderAscending = "_EventsTableSortOrderAscending"
)

const (
	RedisDefaultHost = "localhost"
	RedisDefaultPort = 6379

	DefaultEventsTopic = "edgex/events/device/#"

	DefaultShouldConnectAtStartup        = false
	DefaultEventsTableSortOrderAscending = false
)

type Config struct {
	RedisHost *string
	RedisPort *int

	EventsTopic string

	ShouldConnectAtStartup        bool
	EventsTableSortOrderAscending bool
}

func GetConfig() *Config {
	a := fyne.CurrentApp()

	return &Config{
		RedisHost:                     String(a.Preferences().StringWithFallback(PrefRedisHost, RedisDefaultHost)),
		RedisPort:                     Int(a.Preferences().IntWithFallback(PrefRedisPort, RedisDefaultPort)),
		EventsTopic:                   DefaultEventsTopic,
		ShouldConnectAtStartup:        a.Preferences().BoolWithFallback(PrefShouldConnectAtStartup, DefaultShouldConnectAtStartup),
		EventsTableSortOrderAscending: a.Preferences().BoolWithFallback(PrefEventsTableSortOrderAscending, DefaultEventsTableSortOrderAscending),
	}
}

func DefaultConfig() *Config {
	return &Config{
		RedisHost:                     String(RedisDefaultHost),
		RedisPort:                     Int(RedisDefaultPort),
		EventsTopic:                   DefaultEventsTopic,
		ShouldConnectAtStartup:        DefaultShouldConnectAtStartup,
		EventsTableSortOrderAscending: DefaultEventsTableSortOrderAscending,
	}
}

// String returns a pointer to the given string.
func String(s string) *string {
	return &s
}

// StringVal returns the value of the string at the pointer, or "" if the
// pointer is nil.
func StringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int returns a pointer to the given int.
func Int(i int) *int {
	return &i
}

// IntVal returns the value of the int at the pointer, or 0 if the pointer is
// nil.
func IntVal(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Float returns a pointer to the given float64.
func Float(f float64) *float64 {
	return &f
}

// FloatVal returns the value of the float64 at the pointer, or 0 if the pointer is
// nil.
func FloatVal(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
