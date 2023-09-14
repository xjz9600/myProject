//go:build e2e

package file

import (
	"myProject/web"
	"testing"
)

func TestDownLoader(t *testing.T) {
	server := web.NewServer()
	downloader := NewDownLoader("../testdata/download")
	server.GET("/download", downloader.Build())
	server.Start(":8083")
}
