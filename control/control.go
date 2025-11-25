package control

import (
	"MagNet_Mesh/tunnel"
	"sync"
)

type Controller struct {
	Tunnels []tunnel.Tunnel
	Done    chan struct{}
	//once    sync.Once
	bufpool *sync.Pool
}

func (c *Controller) handlePeer() {
	// Implement peer handling logic here
}

func ListenPeer(p tunnel.Peer, handler func([]byte)) {
	for {
		select {
		case <-p.Done():
			return
		default:
		}
		buffer, err := p.Read()
		go func(buffer []byte) {
			if err != nil {
				return
			}
			handler(buffer)
			p.Put(buffer)
		}(buffer)
	}
}

func handleData(data []byte) {

}

func NewController(bufferSize int) *Controller { //uniform bufpool,default 2048 bytes buffer to keep memory alignment
	if bufferSize <= 0 {
		bufferSize = 2048
	}
	return &Controller{
		Done: make(chan struct{}),
		bufpool: &sync.Pool{
			New: func() any {
				return make([]byte, bufferSize)
			},
		},
	}
}

func (c *Controller) Start() {
	for _, tunnel := range c.Tunnels {
		_ = tunnel.Start()
	}
	for _, tunnel := range c.Tunnels {
		for _, peer := range tunnel.GetPeerList() {
			go ListenPeer(peer, handleData)
		}
	}
}

func (c *Controller) Stop() {
	for _, tunnel := range c.Tunnels {
		_ = tunnel.Stop()
		for _, peer := range tunnel.GetPeerList() {
			_ = peer.Close()
		}
	}
	close(c.Done)
}
