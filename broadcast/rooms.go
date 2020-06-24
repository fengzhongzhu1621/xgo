package broadcast

import "github.com/dustin/go-broadcast"

var roomChannels = make(map[string]broadcast.Broadcaster)

func OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	Room(roomid).Register(listener)
	return listener
}

func CloseListener(roomid string, listener chan interface{}) {
	Room(roomid).Unregister(listener)
	close(listener)
}

func Room(roomid string) broadcast.Broadcaster {
	b, ok := roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		roomChannels[roomid] = b
	}
	return b
}
