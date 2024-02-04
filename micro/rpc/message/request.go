package message

import (
	"bytes"
	"encoding/binary"
)

type Request struct {
	HeadLength uint32
	BodyLength uint32
	RequestID  uint32
	Version    uint8
	Compress   uint8
	Serializer uint8

	ServiceName string
	MethodName  string
	Meta        map[string]string

	Data []byte
}

func (r *Request) CalculateHeaderLength() {
	r.HeadLength = uint32(15 + len(r.ServiceName) + 1 + len(r.MethodName) + 1)
	for k, v := range r.Meta {
		r.HeadLength += uint32(len(k))
		r.HeadLength++
		r.HeadLength += uint32(len(v))
		r.HeadLength++
	}
}

func (r *Request) CalculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}

func EncodeReq(res *Request) []byte {
	resBytes := make([]byte, res.HeadLength+res.BodyLength)
	binary.BigEndian.PutUint32(resBytes[:4], res.HeadLength)
	binary.BigEndian.PutUint32(resBytes[4:8], res.BodyLength)
	binary.BigEndian.PutUint32(resBytes[8:12], res.RequestID)
	resBytes[12] = res.Version
	resBytes[13] = res.Compress
	resBytes[14] = res.Serializer

	headLength := resBytes[15:res.HeadLength]
	copy(headLength[:len(res.ServiceName)], res.ServiceName)
	headLength = headLength[len(res.ServiceName):]
	headLength[0] = '\n'
	headLength = headLength[1:]
	copy(headLength[:len(res.MethodName)], res.MethodName)
	headLength = headLength[len(res.MethodName):]
	headLength[0] = '\n'
	headLength = headLength[1:]
	if len(res.Meta) > 0 {
		for k, v := range res.Meta {
			copy(headLength[:len(k)], k)
			headLength = headLength[len(k):]
			headLength[0] = '\r'
			headLength = headLength[1:]
			copy(headLength[:len(v)], v)
			headLength = headLength[len(v):]
			headLength[0] = '\n'
			headLength = headLength[1:]
		}
	}
	if len(res.Data) > 0 {
		copy(resBytes[res.HeadLength:], res.Data)
	}
	return resBytes
}

func DecodeReq(resBytes []byte) *Request {
	res := &Request{}
	headLength := binary.BigEndian.Uint32(resBytes[:4])
	res.HeadLength = headLength
	bodyLength := binary.BigEndian.Uint32(resBytes[4:8])
	res.BodyLength = bodyLength
	requestID := binary.BigEndian.Uint32(resBytes[8:12])
	res.RequestID = requestID
	res.Version = resBytes[12]
	res.Compress = resBytes[13]
	res.Serializer = resBytes[14]

	headLengthBytes := resBytes[15:headLength]
	idx := bytes.IndexByte(headLengthBytes, '\n')
	res.ServiceName = string(headLengthBytes[:idx])
	headLengthBytes = headLengthBytes[idx+1:]
	idx = bytes.IndexByte(headLengthBytes, '\n')
	res.MethodName = string(headLengthBytes[:idx])
	headLengthBytes = headLengthBytes[idx+1:]
	// 有meta
	idx = bytes.IndexByte(headLengthBytes, '\n')
	if idx != -1 {
		meta := make(map[string]string)
		for idx != -1 {
			metaIndex := bytes.IndexByte(headLengthBytes, '\r')
			key := headLengthBytes[:metaIndex]
			value := headLengthBytes[metaIndex+1 : idx]
			meta[string(key)] = string(value)
			headLengthBytes = headLengthBytes[idx+1:]
			idx = bytes.IndexByte(headLengthBytes, '\n')
		}
		res.Meta = meta
	}
	if uint32(len(resBytes)) > headLength {
		bodyLengthBytes := resBytes[headLength:]
		res.Data = bodyLengthBytes
	}
	return res
}
