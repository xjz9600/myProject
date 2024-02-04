package rpc

import (
	"context"
	"myProject/micro/rpc/compress"
	"myProject/micro/rpc/compress/gz"
	"strconv"
	"time"
	//"encoding/json"
	"errors"
	"myProject/micro/rpc/message"
	"myProject/micro/rpc/serialize"
	"myProject/micro/rpc/serialize/json"
	"net"
	"reflect"
)

type Server struct {
	services  map[string]reflectionStub
	serialize map[uint8]serialize.Serializer
	compress  map[uint8]compress.Compress
}

func NewServer() *Server {
	server := &Server{
		services:  make(map[string]reflectionStub, 16),
		serialize: make(map[uint8]serialize.Serializer, 4),
		compress:  make(map[uint8]compress.Compress, 4),
	}
	server.RegisterSerialize(json.Serializer{})
	server.RegisterCompress(gz.GzipCompress{})
	return server
}
func (s *Server) RegisterSerialize(sl serialize.Serializer) {
	s.serialize[sl.Code()] = sl
}

func (s *Server) RegisterCompress(cp compress.Compress) {
	s.compress[cp.Code()] = cp
}

func (s *Server) Register(service Service) {
	s.services[service.Name()] = reflectionStub{
		s:         service,
		Value:     reflect.ValueOf(service),
		serialize: s.serialize,
		compress:  s.compress,
	}
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
			if err := s.handleConn(conn); err != nil {
				conn.Close()
			}
		}()

	}
}

func (s *Server) handleConn(con net.Conn) error {
	for {
		msg, err := ReadMsg(con)
		if err != nil {
			return err
		}
		req := message.DecodeReq(msg)
		ctx := context.Background()
		if deadlineStr, ok := req.Meta["deadline"]; ok {
			if deadline, er := strconv.ParseInt(deadlineStr, 10, 64); er == nil {
				ctx, _ = context.WithDeadline(ctx, time.UnixMilli(deadline))
			}
		}
		oneway, ok := req.Meta["one-way"]
		if ok && oneway == "true" {
			ctx = CtxWithOneway(ctx)
			go func() {
				s.Invoke(ctx, req)
			}()
			continue
		}
		resp, err := s.Invoke(ctx, req)
		if err != nil {
			resp.Error = []byte(err.Error())
		}
		resp.CalculateHeaderLength()
		resp.CalculateBodyLength()
		res := message.EncodeResp(resp)
		_, err = con.Write(res)
		if err != nil {
			return err
		}
	}
}

func (s *Server) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	res := &message.Response{
		RequestID:  req.RequestID,
		Version:    req.Version,
		Compress:   req.Compress,
		Serializer: req.Serializer,
	}
	service, ok := s.services[req.ServiceName]
	if !ok {
		return res, errors.New("你要调用的服务不存在")
	}
	resp, err := service.invoke(ctx, req)
	res.Data = resp
	return res, err
}

type reflectionStub struct {
	s Service
	reflect.Value
	serialize map[uint8]serialize.Serializer
	compress  map[uint8]compress.Compress
}

func (s *reflectionStub) invoke(ctx context.Context, req *message.Request) ([]byte, error) {
	method := s.Value.MethodByName(req.MethodName)
	inReq := reflect.New(method.Type().In(1).Elem())
	sl, ok := s.serialize[req.Serializer]
	if !ok {
		return nil, errors.New("micro：不支持的协议")
	}
	cp, ok := s.compress[req.Compress]
	if !ok {
		return nil, errors.New("micro：不支持的压缩算法")
	}
	data, err := cp.UnCompression(req.Data)
	if err != nil {
		return nil, err
	}
	err = sl.Decode(data, inReq.Interface())
	if err != nil {
		return nil, err
	}

	serResp := method.Call([]reflect.Value{reflect.ValueOf(ctx), inReq})
	if serResp[1].Interface() != nil {
		err = serResp[1].Interface().(error)
	}
	if serResp[0].IsNil() {
		return nil, err
	}
	res, er := sl.Encode(serResp[0].Interface())
	if er != nil {
		return nil, er
	}
	data, er = cp.Compression(res)
	if er != nil {
		return nil, er
	}
	return data, err
}
