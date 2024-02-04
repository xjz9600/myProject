package micro

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

type Pool struct {
	idlesCons   chan *idleConn
	maxCnt      int
	maxIdleCnt  int
	currentCnt  int
	waitQueue   []chan *queueReq
	maxIdleTime time.Duration
	factory     func() (net.Conn, error)
	mutex       sync.RWMutex
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	// 进来先判断超时
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for {
		select {
		case idle := <-p.idlesCons:
			// 空闲链接超过最大存活时间
			if idle.timeout.Add(p.maxIdleTime).Before(time.Now()) {
				p.mutex.Lock()
				p.currentCnt--
				p.mutex.Unlock()
				_ = idle.c.Close()
				continue
			}
			return idle.c, nil
		default:
			p.mutex.Lock()
			if p.currentCnt >= p.maxCnt {
				req := make(chan *queueReq, 1)
				p.waitQueue = append(p.waitQueue, req)
				p.mutex.Unlock()
				select {
				case <-ctx.Done():
					go func() {
						// 归还返回的链接
						c := <-req
						p.Put(context.Background(), c.c)
					}()
					return nil, ctx.Err()
				case c := <-req:
					return c.c, nil
				}
			}
			newConn, err := p.factory()
			if err != nil {
				return nil, err
			}
			p.currentCnt++
			p.mutex.Unlock()
			return newConn, nil
		}
	}
}

func (p *Pool) Put(ctx context.Context, c net.Conn) error {
	// 进来先判断超时
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	p.mutex.Lock()
	if len(p.waitQueue) > 0 {
		wq := p.waitQueue[0]
		p.waitQueue = p.waitQueue[1:]
		wq <- &queueReq{c: c}
		p.mutex.Unlock()
		return nil
	}
	defer p.mutex.Unlock()
	idle := &idleConn{c: c, timeout: time.Now()}
	select {
	case p.idlesCons <- idle:
	default:
		p.currentCnt--
		_ = idle.c.Close()
	}
	return nil
}

func NewPool(maxCnt, initCnt, maxIdleCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, errors.New("micro：初始化链接不能超过最大空闲链接数")
	}
	idlesCons := make(chan *idleConn, initCnt)
	for i := 0; i < initCnt; i++ {
		c, err := factory()
		if err != nil {
			return nil, err
		}
		ic := &idleConn{
			c:       c,
			timeout: time.Now(),
		}
		idlesCons <- ic
	}
	pool := &Pool{
		maxCnt:      maxCnt,
		maxIdleCnt:  maxIdleCnt,
		factory:     factory,
		maxIdleTime: maxIdleTime,
	}
	return pool, nil
}

type queueReq struct {
	c net.Conn
}

type idleConn struct {
	c       net.Conn
	timeout time.Time
}
