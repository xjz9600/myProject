package channel

import (
	"errors"
	"net"
	"time"
)

type Pool struct {
	idlesCons   chan *idleConn
	maxCnt      int
	maxIdleCnt  int
	currentCnt  int
	waitQueue   chan []*queueReq
	maxIdleTime time.Duration
	factory     func(conn net.Conn, err error)
}

func NewPool(maxCnt, initCnt, maxIdleCnt int, maxIdleTime time.Duration, factory func(conn net.Conn, err error)) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, errors.New("")
	}
}

type queueReq struct {
	c net.Conn
}

type idleConn struct {
	c       net.Conn
	timeout time.Duration
}
