package rpc

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"myProject/micro/rpc/compress/gz"
	"myProject/micro/rpc/message"
	"myProject/micro/rpc/serialize/json"
	"testing"
)

func Test_setFuncField(t *testing.T) {
	ctrl := gomock.NewController(t)
	testCases := []struct {
		name         string
		mockService  Service
		mockProxy    func() Proxy
		wantErr      error
		wantProxyErr error
	}{
		{
			name:        "service nil",
			mockService: nil,
			mockProxy: func() Proxy {
				return NewMockProxy(ctrl)
			},
			wantErr: errors.New("rpc：不支持nil"),
		},
		{
			name:        "get response",
			mockService: &UserService{},
			mockProxy: func() Proxy {
				mock := NewMockProxy(ctrl)
				req := &message.Request{
					HeadLength:  36,
					BodyLength:  10,
					Serializer:  1,
					MethodName:  "GetById",
					ServiceName: "user_service",
					Data:        []byte(`{"Id":123}`),
				}
				mock.EXPECT().Invoke(gomock.Any(), req).Return(&message.Response{}, nil)
				return mock
			},
		},
		{
			name:        "get response error",
			mockService: &UserService{},
			mockProxy: func() Proxy {
				mock := NewMockProxy(ctrl)
				mock.EXPECT().Invoke(gomock.Any(), &message.Request{
					HeadLength:  36,
					BodyLength:  10,
					Serializer:  1,
					MethodName:  "GetById",
					ServiceName: "user_service",
					Data:        []byte(`{"Id":123}`),
				}).Return(&message.Response{}, errors.New("response err"))
				return mock
			},
			wantProxyErr: errors.New("response err"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := setFuncField(tc.mockService, tc.mockProxy(), json.Serializer{}, gz.GzipCompress{})
			assert.Equal(t, err, tc.wantErr)
			if err != nil {
				return
			}
			_, err = tc.mockService.(*UserService).GetById(context.Background(), &GetByIdReq{
				Id: 123,
			})
			assert.Equal(t, err, tc.wantProxyErr)
		})
	}
}
