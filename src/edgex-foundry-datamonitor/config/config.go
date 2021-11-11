package config

import "fyne.io/fyne/v2"

type Config struct {
	app         fyne.App
	EventsTopic string
}

func GetConfig() *Config {
	a := fyne.CurrentApp()

	return &Config{
		app:         a,
		EventsTopic: DefaultEventsTopic,
	}
}

func (c *Config) GetRedisHost() string {
	return c.app.Preferences().StringWithFallback(PrefRedisHost, RedisDefaultHost)
}

func (c *Config) GetRedisPort() int {
	return c.app.Preferences().IntWithFallback(PrefRedisPort, RedisDefaultPort)
}

func (c *Config) GetShouldConnectAtStartup() bool {
	return c.app.Preferences().BoolWithFallback(PrefShouldConnectAtStartup, DefaultShouldConnectAtStartup)
}

func (c *Config) GetEventsTableSortOrderAscending() bool {
	return c.app.Preferences().BoolWithFallback(PrefEventsTableSortOrderAscending, DefaultEventsTableSortOrderAscending)
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
