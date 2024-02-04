package micro

import (
	"encoding/binary"
	"net"
	"time"
)

type Client struct {
	network string
	addr    string
}

func NewClient(network, add string) *Client {
	return &Client{
		network: network,
		addr:    add,
	}
}

func (c *Client) Send(msg string) ([]byte, error) {
	conn, err := net.DialTimeout(c.network, c.addr, time.Second*3)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	data := make([]byte, len(msg)+msgLengthBytes)
	binary.BigEndian.PutUint64(data[:msgLengthBytes], uint64(len(msg)))
	copy(data[msgLengthBytes:], msg)
	_, err = conn.Write(data)
	if err != nil {
		return nil, err
	}
	lenBytes := make([]byte, msgLengthBytes)
	_, err = conn.Read(lenBytes)
	if err != nil {
		return nil, err
	}
	msgLen := binary.BigEndian.Uint64(lenBytes)
	msgInfo := make([]byte, msgLen)
	_, err = conn.Read(msgInfo)
	if err != nil {
		return nil, err
	}
	return msgInfo, nil
}
