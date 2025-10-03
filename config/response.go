package config

// EventType defines the event type of config change.
type EventType uint8

const (
	// EventTypeNull represents null event.
	EventTypeNull EventType = 0

	// EventTypePut represents set or update config event.
	EventTypePut EventType = 1

	// EventTypeDel represents delete config event.
	EventTypeDel EventType = 2
)

// Response defines config center's response interface.
type IResponse interface {
	// Value returns config value as string.
	Value() string

	// MetaData returns extra metadata. With option,
	// we can implement some extra features for different config center,
	// such as namespace, group, lease, etc.
	MetaData() map[string]string

	// Event returns the type of watch event.
	Event() EventType
}
