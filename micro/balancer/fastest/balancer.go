package fastest

import (
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"net/http"
	"net/url"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

type BalancerBuilder struct {
	Point    string
	Query    string
	Duration time.Duration
}

func (b *BalancerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	con := make([]*conn, 0, len(info.ReadySCs))
	for k, v := range info.ReadySCs {
		wc := &conn{
			c:        k,
			addr:     v.Address,
			response: 100 * time.Millisecond,
		}
		con = append(con, wc)
	}
	bal := &Balancer{connects: con}
	closeChan := make(chan struct{})
	runtime.SetFinalizer(bal, func(b *Balancer) {
		closeChan <- struct{}{}
	})
	go func() {
		ticker := time.NewTicker(b.Duration)
		select {
		case <-ticker.C:
			bal.updateRespTime(b.Point, b.Query)
		case <-closeChan:
			return
		}
	}()
	return bal
}

func (b *Balancer) updateRespTime(endpoint, query string) {
	info := &QueryInfo{}
	uStr := endpoint + "/api/v1/query?query=" + query
	u, err := url.Parse(uStr)
	if err != nil {
		return
	}
	u.RawQuery = u.Query().Encode()
	err = GetPromResult(u.String(), &info)
	if err != nil {
		return
	}
	for _, in := range info.Data.Result {
		promAddr := in.Metric["addr"]
		for _, c := range b.connects {
			if c.addr.Addr == promAddr {
				ms, er := strconv.ParseInt(in.Value[1].(string), 10, 64)
				if er != nil {
					continue
				}
				b.mutex.Lock()
				c.response = time.Duration(ms) * time.Millisecond
				b.mutex.Unlock()
			}
		}
	}
}

type Balancer struct {
	connects []*conn
	mutex    sync.RWMutex
}

func (b *Balancer) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	b.mutex.RLock()
	if len(b.connects) == 0 {
		b.mutex.RUnlock()
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var conn *conn
	for _, cn := range b.connects {
		if conn == nil || cn.response < conn.response {
			conn = cn
		}
	}
	b.mutex.RUnlock()
	return balancer.PickResult{
		SubConn: conn.c,
		Done: func(info balancer.DoneInfo) {
		},
	}, nil
}

type conn struct {
	c        balancer.SubConn
	addr     resolver.Address
	response time.Duration
}

type ResultType struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value"`
}

type QueryData struct {
	ResultType string       `json:"resultType"`
	Result     []ResultType `json:"result"`
}

type QueryInfo struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

func GetPromResult(url string, result interface{}) error {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(result)
	if err != nil {
		fmt.Printf("%s", debug.Stack())
		debug.PrintStack()
		return err
	}
	return nil
}
