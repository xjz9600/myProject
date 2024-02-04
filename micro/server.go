package micro

import (
	"encoding/binary"
	"net"
)

const msgLengthBytes = 8

type Server struct {
}

func (s *Server) Serve(network, addr string) error {
	listener, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go func() {
			if err := s.sendMsg(conn); err != nil {
				conn.Close()
			}
		}()

	}
}

func (s *Server) sendMsg(con net.Conn) error {
	lenBytes := make([]byte, msgLengthBytes)
	_, err := con.Read(lenBytes)
	if err != nil {
		return err
	}
	msgLen := binary.BigEndian.Uint64(lenBytes)
	msg := make([]byte, msgLen)
	_, err = con.Read(msg)
	if err != nil {
		return err
	}
	data := genMsg(msg)
	res := make([]byte, len(data)+msgLengthBytes)
	binary.BigEndian.PutUint64(res[:msgLengthBytes], uint64(len(data)))
	copy(res[msgLengthBytes:], data)
	_, err = con.Write(res)
	if err != nil {
		return err
	}
	return nil
}

func genMsg(msg []byte) []byte {
	data := make([]byte, len(msg)*2)
	copy(data[:len(msg)], msg)
	copy(data[len(msg):], msg)
	return data
}
