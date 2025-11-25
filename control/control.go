package control

import (
	"MagNet_Mesh/tunnel"
	"sync"

	"github.com/panjf2000/ants"
)

type Controller struct {
	Tunnels    []tunnel.Tunnel
	Done       chan struct{}
	BufPool    *sync.Pool
	In         chan []byte //channel to receive data from tunnels
	WorkerPool *ants.Pool
}

func NewController(bufferSize int, workerSize int, channelSize int) *Controller {
	if bufferSize <= 0 {
		bufferSize = 2048
	}
	if workerSize <= 0 {
		workerSize = 1000
	}
	if channelSize <= 0 {
		channelSize = 64
	}
	workerPool, _ := ants.NewPool(workerSize) //default 1000 workers in the pool
	return &Controller{
		Done: make(chan struct{}),
		BufPool: &sync.Pool{ //uniform buffer pool,default 2048 bytes buffer to keep memory alignment
			New: func() any {
				return make([]byte, bufferSize)
			},
		},
		In:         make(chan []byte, channelSize), //buffered channel to receive data from tunnels
		WorkerPool: workerPool,
	}
}

func (c *Controller) Start() {
	for _, tunnel := range c.Tunnels {
		_ = tunnel.Start()
	}
}

func (c *Controller) Stop() {
	for _, tunnel := range c.Tunnels {
		_ = tunnel.Stop()
	}
	close(c.Done)
}

func (c *Controller) Run() {
	for {
		select {
		case data := <-c.In:
			_ = c.WorkerPool.Submit(func() {
				handleData(data)
				c.BufPool.Put(data)
			})
		case <-c.Done:
			return
		}
	}
}

func handleData(data []byte) {
	// Implement data handling logic here

}
