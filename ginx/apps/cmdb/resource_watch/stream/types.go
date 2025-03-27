package stream

type OperType string

// OperType TODO

const (
	// Insert TODO
	// reference doc:
	// https://docs.mongodb.com/manual/reference/change-events/#change-events
	// Document operation type
	Insert OperType = "insert"
	// Delete TODO
	Delete OperType = "delete"
	// Replace TODO
	Replace OperType = "replace"
	// Update TODO
	Update OperType = "update"

	// Drop TODO
	// collection operation type.
	Drop OperType = "drop"
	// Rename TODO
	Rename OperType = "rename"

	// DropDatabase event occurs when a database is dropped.
	DropDatabase OperType = "dropDatabase"

	// Invalidate TODO
	// For change streams opened up against a collection, a drop event, rename event,
	// or dropDatabase event that affects the watched collection leads to an invalidate event.
	Invalidate OperType = "invalidate"

	// Lister OperType is a self defined type, which is represent this operation comes from
	// a list watcher's find operations, it does not really come form the mongodb's change event.
	Lister OperType = "lister"
	// ListDone OperType is a self defined type, which means that the list operation has already finished,
	// and the watch events starts. this OperType send only for once.
	// Note: it's only used in the ListWatch Operation.
	ListDone OperType = "listerDone"
)
