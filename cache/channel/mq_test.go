package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroker_Send(t *testing.T) {

	b := &Broker{}
	go func() {
		for {
			b.Send(Msg{content: time.Now().String()})
		}
	}()
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			msg, _ := b.subscribe(100)
			for m := range msg {
				fmt.Println(m.content)
			}
		}()
	}
	wg.Wait()
}
