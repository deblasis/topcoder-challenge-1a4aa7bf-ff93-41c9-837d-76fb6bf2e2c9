package services

type ConnectionState int

const (
	ClientDisconnected ConnectionState = iota
	ClientConnecting
	ClientConnected
)

type processorState int

const (
	Stopped processorState = iota
	Paused
	Running
)
