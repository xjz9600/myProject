package micro

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClientAndServer(t *testing.T) {
	go func() {
		server := Server{}
		server.Serve("tcp", "localhost:8082")
	}()
	time.Sleep(time.Second * 3)
	client := NewClient("tcp", "localhost:8082")
	for i := 0; i < 5; i++ {
		req, err := client.Send("Hello Word")
		assert.NoError(t, err)
		assert.Equal(t, string(req), "Hello WordHello Word")
	}
}
