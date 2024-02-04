package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeAndDecodeResp(t *testing.T) {
	testCases := []struct {
		name    string
		wantReq *Response
	}{
		{
			name: "normal",
			wantReq: &Response{
				RequestID:  123,
				Version:    3,
				Compress:   2,
				Serializer: 5,
				Error:      []byte("test err"),
				Data:       []byte("req data info"),
			},
		},
		{
			name: "no err",
			wantReq: &Response{
				RequestID:  123,
				Version:    3,
				Compress:   2,
				Serializer: 5,
				Data:       []byte("req data info"),
			},
		},
		{
			name: "no data",
			wantReq: &Response{
				RequestID:  123,
				Version:    3,
				Compress:   2,
				Serializer: 5,
				Error:      []byte("test err"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.wantReq.CalculateBodyLength()
			tc.wantReq.CalculateHeaderLength()
			reqBytes := EncodeResp(tc.wantReq)
			req := DecodeResp(reqBytes)
			assert.Equal(t, tc.wantReq, req)
		})
	}
}
