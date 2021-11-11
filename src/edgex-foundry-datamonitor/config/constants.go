package config

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
