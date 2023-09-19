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
	conns    *circularQueue
	mu       sync.Mutex
	IdleTime time.Duration // 閒置時間
}

type connWithTime struct {
	net.Conn
	t time.Time
}

func (p *ConnPool) CreatePool() {
	p.conns = newCircularQueue(p.MaxIdle)

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

		p.conns.enqueue(con)
	}
}

func (p *ConnPool) GetConnsLen() int {
	return p.conns.size()
}

func (p *ConnPool) getInstance() (net.Conn, error) {
	conn := p.conns.dequeue()
	if conn == nil {
		conn, err := p.Dial()
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	if p.IdleTime > 0 && time.Since(conn.t) > p.IdleTime { // 判断闲置并移除
		conn.Close()
		conn = nil
		return p.getInstance()
	}

	return conn.Conn, nil
}

func (p *ConnPool) Get() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.getInstance()
}

func (p *ConnPool) Put(conn net.Conn) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.conns.enqueue(&connWithTime{conn, time.Now()})
}

func (p *ConnPool) Release() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.conns.each(func(node *connWithTime) {
		node.Close()
	})
	return p.conns.clear()
}
