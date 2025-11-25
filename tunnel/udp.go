package tunnel

import (
	"errors"
	"net"
	"sync"
)

type UdpTunnel struct {
	localAddr *net.UDPAddr //UDP local listen address

	conn    *net.UDPConn //UDP instance
	closed  chan struct{}
	BufPool *sync.Pool //Buffer pool for reusing byte slices
	in      chan<- []byte
}

func (u *UdpTunnel) Done() <-chan struct{} {
	return u.closed
}

func (u *UdpTunnel) Type() string {
	return "UDPTunnel"
}

func NewUdpTunnel(localAddr *net.UDPAddr, bufPool *sync.Pool, in chan<- []byte) *UdpTunnel {
	return &UdpTunnel{
		localAddr: localAddr,
		closed:    make(chan struct{}),
		BufPool:   bufPool,
		in:        in,
	}
}

func (u *UdpTunnel) Start() error {
	if u.conn == nil {
		if u.localAddr == nil {
			return errors.New("no local address")
		}
		conn, err := net.ListenUDP("udp", u.localAddr)
		if err != nil {
			return err
		}
		u.conn = conn
	}
	for {
		select {
		case <-u.closed:
			break
		default:
		}
		buffer := u.BufPool.Get().([]byte)
		n, addr, err := u.conn.ReadFromUDP(buffer)
		if err != nil {
			u.BufPool.Put(buffer)
			continue
		}
		// Process data from addr
		go func(buffer []byte, n int, addr *net.UDPAddr) {
			select {
			case u.in <- buffer:
			default:
				u.BufPool.Put(buffer)
			}
		}(buffer, n, addr)
	}
}

func (u *UdpTunnel) Stop() error {
	close(u.closed)
	return u.conn.Close()
}

func (u *UdpTunnel) Put(data []byte) {
	u.BufPool.Put(data)
}

//func (u *UdpTunnel) GetPeerList() []Peer {
//	u.mu.RLock()
//	defer u.mu.RUnlock()
//
//	// 创建一个新的 map 用于返回
//	peers := make([]Peer, 0, len(u.peers))
//	for _, v := range u.peers {
//		peers = append(peers, v)
//	}
//	return peers
//}
