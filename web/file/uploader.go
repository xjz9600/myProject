package file

import (
	"io"
	"mime/multipart"
	"myProject/web"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type fileUploader struct {
	fileField  string
	dstPathFun func(*multipart.FileHeader) string
}

type fileUploaderOpt func(*fileUploader)

func NewFileUploader(opts ...fileUploaderOpt) *fileUploader {
	res := &fileUploader{
		fileField: "myFile",
		dstPathFun: func(fileHeader *multipart.FileHeader) string {
			return filepath.Join("../testdata", "upload", fileHeader.Filename)
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func WithFileField(fileField string) fileUploaderOpt {
	return func(uploader *fileUploader) {
		uploader.fileField = fileField
	}
}

func WithDstPathFun(dstPathFun func(*multipart.FileHeader) string) fileUploaderOpt {
	return func(uploader *fileUploader) {
		uploader.dstPathFun = dstPathFun
	}
}

func (f *fileUploader) Handle(ctx *web.Context) {
	filedata, fileHeader, err := ctx.Req.FormFile(f.fileField)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	dst := f.dstPathFun(fileHeader)
	dir, _ := path.Split(dst)
	os.MkdirAll(dir, os.ModePerm)
	file, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	defer file.Close()
	_, err = io.CopyBuffer(file, filedata, nil)
	if err != nil {
		ctx.RespStatusCode = http.StatusInternalServerError
		ctx.RespData = []byte(err.Error())
		return
	}
	ctx.RespStatusCode = http.StatusOK
	ctx.RespData = []byte("上传成功")
}
