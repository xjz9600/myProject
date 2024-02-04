package rpc

import (
	"encoding/binary"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	lenBytes := make([]byte, msgLengthBytes)
	_, err := conn.Read(lenBytes)
	if err != nil {
		return nil, err
	}
	headLength := binary.BigEndian.Uint32(lenBytes[:4])
	bodyLength := binary.BigEndian.Uint32(lenBytes[4:])
	msg := make([]byte, headLength+bodyLength)
	copy(msg[:8], lenBytes)
	_, err = conn.Read(msg[8:])
	if err != nil {
		return nil, err
	}
	return msg, nil
}
