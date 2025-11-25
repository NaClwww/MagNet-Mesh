package tunnel

import (
	"errors"
	"net"
	"sync"
	"time"
)

type UdpPeer struct {
	addr     *net.UDPAddr
	lastSeen time.Time

	tunnel *UdpTunnel
	in     chan []byte
	done   chan struct{}
	once   sync.Once
}

func NewUdpPeer(addr *net.UDPAddr) *UdpPeer {
	return &UdpPeer{
		addr: addr,
		in:   make(chan []byte, 16), // 可配置
		done: make(chan struct{}),
	}
}

func (u *UdpPeer) Done() <-chan struct{} {
	return u.done
}

func (u *UdpPeer) Read() ([]byte, error) {
	select {
	case <-u.done:
		return nil, ErrClosed
	case b, ok := <-u.in:
		if !ok {
			return nil, ErrClosed
		}
		u.lastSeen = time.Now()
		return b, nil
	}
}

func (u *UdpPeer) Put(b []byte) {
	if u.tunnel != nil {
		u.tunnel.bufpool.Put(b)
		return
	}
	// 没有 tunnel 时安全丢弃
}

func (u *UdpPeer) Send(data []byte) (int, error) {
	return u.tunnel.conn.WriteToUDP(data, u.addr)
}

func (u *UdpPeer) Close() error {
	u.once.Do(func() { close(u.done) })
	return nil
}

func (u *UdpPeer) LastSeen() time.Time {
	return u.lastSeen
}

func (u *UdpPeer) Type() Tunnel {
	return u.tunnel
}

type UdpTunnel struct {
	BufSize int

	localAddr *net.UDPAddr

	conn    *net.UDPConn
	closed  chan struct{}
	bufpool *sync.Pool

	mu    sync.RWMutex
	peers map[string]*UdpPeer
}

func NewUdpTunnel(conn *net.UDPConn, bufpool *sync.Pool) *UdpTunnel {
	return &UdpTunnel{
		BufSize: 2048,
		conn:    conn,
		closed:  make(chan struct{}),
		bufpool: bufpool,
		peers:   make(map[string]*UdpPeer),
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
		buffer := u.bufpool.Get().([]byte)
		n, addr, err := u.conn.ReadFromUDP(buffer)
		if err != nil {
			u.bufpool.Put(buffer)
			continue
		}
		// Process data from addr
		go func(buffer []byte, n int, addr *net.UDPAddr) {
			peer, has := u.peers[addr.String()]
			if !has {
				//TODO: add peer management
				//auth here if needed
				//if auth passed,add peer

				u.bufpool.Put(buffer)
			}
			select {
			case peer.in <- buffer:
			default:
				u.bufpool.Put(buffer)
			}
		}(buffer, n, addr)
	}
}

func (u *UdpTunnel) Stop() error {
	close(u.closed)
	return u.conn.Close()
}

func (u *UdpTunnel) GetPeerList() []Peer {
	u.mu.RLock()
	defer u.mu.RUnlock()

	// 创建一个新的 map 用于返回
	peers := make([]Peer, 0, len(u.peers))
	for _, v := range u.peers {
		peers = append(peers, v)
	}
	return peers
}
