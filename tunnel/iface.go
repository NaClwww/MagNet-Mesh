package tunnel

// for tunnel interface:
// done chan struct{} : to signal tunnel is closed
// in chan []byte : to receive data from tunnel peers,

type Tunnel interface { // Represents a Tunnel that can connect to Peers,Receive data and distribute data
	Start() error
	Stop() error
	Done() <-chan struct{}
	Type() string
	//GetPeerList() []Peer
}

//type Peer interface { // Represents a connected peer over a Tunnel,Sending data
//	Read() ([]byte, error)
//	Send(data []byte) (int, error)
//	Close() error
//	Put(data []byte)
//	LastSeen() time.Time
//	Type() Tunnel
//	Done() <-chan struct{}
//}
