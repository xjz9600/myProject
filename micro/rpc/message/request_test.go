package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecodeReq(t *testing.T) {
	testCases := []struct {
		name    string
		wantReq *Request
	}{
		{
			name: "normal",
			wantReq: &Request{
				RequestID:   123,
				Version:     3,
				Compress:    2,
				Serializer:  5,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"test": "hello word",
				},
				Data: []byte("req\ndata\ninfo"),
			},
		},
		{
			name: "no Meta",
			wantReq: &Request{
				RequestID:   123,
				Version:     3,
				Compress:    2,
				Serializer:  5,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Data:        []byte("req data info"),
			},
		},
		{
			name: "no data",
			wantReq: &Request{
				RequestID:   123,
				Version:     3,
				Compress:    2,
				Serializer:  5,
				ServiceName: "user_service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"test": "hello word",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.wantReq.CalculateBodyLength()
			tc.wantReq.CalculateHeaderLength()
			reqBytes := EncodeReq(tc.wantReq)
			req := DecodeReq(reqBytes)
			assert.Equal(t, tc.wantReq, req)
		})
	}
}
