package tunnel

import (
	"errors"
	"time"
)

var ErrClosed = errors.New("peer closed")

type Tunnel interface { // Represents a Tunnel that can connect to Peers,Receive data and distribute data
	Start() error
	Stop() error
	GetPeerList() []Peer
}

type Peer interface { // Represents a connected peer over a Tunnel,Sending data
	Read() ([]byte, error)
	Send(data []byte) (int, error)
	Close() error
	Put(data []byte)
	LastSeen() time.Time
	Type() Tunnel
	Done() <-chan struct{}
}
