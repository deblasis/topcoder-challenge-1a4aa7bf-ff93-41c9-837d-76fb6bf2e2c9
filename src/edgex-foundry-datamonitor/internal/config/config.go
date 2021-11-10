package config

import "fyne.io/fyne/v2"

const (
	prefRedisHost = "_RedisHost"
	prefRedisPort = "_RedisPort"

	PrefShouldConnectAtStartup = "_ShouldConnectAtStartup"
)

const (
	RedisDefaultHost = "localhost"
	RedisDefaultPort = 6379
)

type Config struct {
	RedisHost *string
	RedisPort *int
}

func GetConfig() *Config {
	a := fyne.CurrentApp()

	return &Config{
		RedisHost: String(a.Preferences().StringWithFallback(prefRedisHost, RedisDefaultHost)),
		RedisPort: Int(a.Preferences().IntWithFallback(prefRedisPort, RedisDefaultPort)),
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
