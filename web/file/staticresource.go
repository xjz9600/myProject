package file

import (
	lru "github.com/hashicorp/golang-lru/v2"
	"myProject/web"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type staticResource struct {
	dst                     string
	extensionContentTypeMap map[string]string
	maxSize                 int
	cache                   *lru.Cache[string, any]
}

type staticResourceOption func(*staticResource)

func NewStaticResource(opts ...staticResourceOption) *staticResource {
	cache, _ := lru.New[string, any](100)
	res := &staticResource{
		dst:     filepath.Join("../testdata/static"),
		maxSize: 1000,
		cache:   cache,
		extensionContentTypeMap: map[string]string{
			// 这里根据自己的需要不断添加
			"jpeg": "image/jpeg",
			"jpe":  "image/jpeg",
			"jpg":  "image/jpeg",
			"png":  "image/png",
			"pdf":  "application/pdf",
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func WithStaticResourceDst(dst string) staticResourceOption {
	return func(resource *staticResource) {
		resource.dst = dst
	}
}

func WithStaticResourceMaxSize(maxSize int) staticResourceOption {
	return func(resource *staticResource) {
		resource.maxSize = maxSize
	}
}

func WithStaticResourceCache(cache *lru.Cache[string, any]) staticResourceOption {
	return func(resource *staticResource) {
		resource.cache = cache
	}
}

func AddStaticResourceExtensionContentTypeMap(typeMap map[string]string) staticResourceOption {
	return func(resource *staticResource) {
		for k, v := range typeMap {
			resource.extensionContentTypeMap[k] = v
		}
	}
}

func (s *staticResource) Handle(ctx *web.Context) {
	staticName, err := ctx.ParamValue("staticName").AsString()
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	path := filepath.Join(s.dst, staticName)
	path, err = filepath.Abs(path)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	header := ctx.Response.Header()
	if contentType, ok := s.extensionContentTypeMap[filepath.Ext(path)[1:]]; ok {
		header.Set("Content-Type", contentType)
	}
	if cacheData, ok := s.cache.Get(path); ok {
		data, _ := cacheData.([]byte)
		header.Set("Content-Length", strconv.Itoa(len(data)))
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = data
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	if len(data) <= s.maxSize {
		s.cache.Add(path, data)
	}
	header.Set("Content-Length", strconv.Itoa(len(data)))
	ctx.RespStatusCode = http.StatusOK
	ctx.RespData = data
	return
}
