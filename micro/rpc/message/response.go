package message

import (
	"encoding/binary"
)

type Response struct {
	HeadLength uint32
	BodyLength uint32
	RequestID  uint32
	Version    uint8
	Compress   uint8
	Serializer uint8

	Error []byte

	Data []byte
}

func EncodeResp(resp *Response) []byte {
	resBytes := make([]byte, resp.HeadLength+resp.BodyLength)
	binary.BigEndian.PutUint32(resBytes[:4], resp.HeadLength)
	binary.BigEndian.PutUint32(resBytes[4:8], resp.BodyLength)
	binary.BigEndian.PutUint32(resBytes[8:12], resp.RequestID)
	resBytes[12] = resp.Version
	resBytes[13] = resp.Compress
	resBytes[14] = resp.Serializer
	if len(resp.Error) > 0 {
		copy(resBytes[15:resp.HeadLength], resp.Error)
	}
	if len(resp.Data) > 0 {
		copy(resBytes[resp.HeadLength:], resp.Data)
	}
	return resBytes
}

func DecodeResp(respBytes []byte) *Response {
	res := &Response{}
	headLength := binary.BigEndian.Uint32(respBytes[:4])
	res.HeadLength = headLength
	bodyLength := binary.BigEndian.Uint32(respBytes[4:8])
	res.BodyLength = bodyLength
	requestID := binary.BigEndian.Uint32(respBytes[8:12])
	res.RequestID = requestID
	res.Version = respBytes[12]
	res.Compress = respBytes[13]
	res.Serializer = respBytes[14]
	if headLength > 15 {
		errBytes := respBytes[15:headLength]
		res.Error = errBytes
	}
	if uint32(len(respBytes)) > headLength {
		bodyLengthBytes := respBytes[headLength:]
		res.Data = bodyLengthBytes
	}
	return res
}

func (r *Response) CalculateHeaderLength() {
	r.HeadLength = uint32(15 + len(r.Error))
}

func (r *Response) CalculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}
