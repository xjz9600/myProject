package proto

import (
	"errors"
	"github.com/xjz9600/protobuf/proto"
)

type Serializer struct {
}

func (s Serializer) Code() uint8 {
	return 2
}

func (s Serializer) Encode(val any) ([]byte, error) {
	message, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("micro: 必须是 proto.Message")
	}
	return proto.Marshal(message)
}

func (s Serializer) Decode(data []byte, val any) error {
	message, ok := val.(proto.Message)
	if !ok {
		return errors.New("micro: 必须是 proto.Message")
	}
	return proto.Unmarshal(data, message)
}
