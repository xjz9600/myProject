//go:build e2e

package file

import (
	"myProject/web"
	"testing"
)

func TestStaticResource(t *testing.T) {
	server := web.NewServer()
	staticResource := NewStaticResource()
	server.GET("/static/:staticName", staticResource.Handle)
	server.Start(":8083")
}
