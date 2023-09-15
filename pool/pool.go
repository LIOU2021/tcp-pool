package pool

import (
	"log"
	"net"
	"sync"
	"time"
)

type ConnPool struct {
	Dial     func() (net.Conn, error) // 連線方式
	MaxIdle  int                      // 最大閒置連線數量
	MinIdle  int                      // 最小閒置連線數量
	conns    chan *connWithTime
	mu       sync.Mutex
	IdleTime time.Duration // 閒置時間
}

type connWithTime struct {
	net.Conn
	t time.Time
}

func (p *ConnPool) CreatePool() {
	p.conns = make(chan *connWithTime, p.MaxIdle)
	for i := 0; i < p.MinIdle; i++ {
		conn, err := p.Dial()
		if err != nil {
			log.Fatal(err)
			return
		}

		con := &connWithTime{
			Conn: conn,
			t:    time.Now(),
		}

		p.conns <- con
	}
}

func (p *ConnPool) GetConnsLen() int {
	return len(p.conns)
}

func (p *ConnPool) Get() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.conns) < 1 {
		conn, err := p.Dial()
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	conn := <-p.conns
	if p.IdleTime > 0 && time.Since(conn.t) > p.IdleTime {
		conn.Close()
		return p.Get()
	}

	return conn.Conn, nil
}

func (p *ConnPool) Put(conn net.Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.MaxIdle > 0 && len(p.conns) >= p.MaxIdle {
		conn.Close()
		return nil
	}
	p.conns <- &connWithTime{conn, time.Now()}
	return nil
}

func (p *ConnPool) Release() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	close(p.conns)
	for i := range p.conns {
		err = i.Close()
		if err != nil {
			return
		}
	}
	return
}
