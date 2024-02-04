package rpc

import (
	"context"
	"errors"
	"github.com/silenceper/pool"
	"myProject/micro/rpc/compress"
	"myProject/micro/rpc/compress/gz"
	"myProject/micro/rpc/message"
	"myProject/micro/rpc/serialize"
	"myProject/micro/rpc/serialize/json"
	"net"
	"reflect"
	"strconv"
	"time"
)

const msgLengthBytes = 8

type ClientOpt func(client *Client)

func (c *Client) InitService(service Service) error {
	// 在这里初始化一个 Proxy
	return setFuncField(service, c, c.serializer, c.compress)
}

func setFuncField(service Service, p Proxy, serializer serialize.Serializer, compress compress.Compress) error {
	if service == nil {
		return errors.New("rpc：不支持nil")
	}
	val := reflect.ValueOf(service)
	if val.Kind() != reflect.Pointer || val.Elem().Kind() != reflect.Struct {
		return errors.New("rpc：只支持一级指针")
	}
	val = val.Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		if fieldVal.CanSet() {
			retVal := reflect.New(fieldVal.Type().Out(0).Elem())
			fn := reflect.MakeFunc(fieldTyp.Type, func(args []reflect.Value) (results []reflect.Value) {
				ctx := args[0].Interface().(context.Context)
				data, err := serializer.Encode(args[1].Interface())
				if err != nil {
					return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
				}
				compressData, err := compress.Compression(data)
				if err != nil {
					return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
				}
				req := &message.Request{
					ServiceName: service.Name(),
					MethodName:  fieldTyp.Name,
					Data:        compressData,
					Serializer:  serializer.Code(),
					Compress:    compress.Code(),
				}
				var meta map[string]string
				if deadline, ok := ctx.Deadline(); ok {
					meta = map[string]string{"deadline": strconv.FormatInt(deadline.UnixMicro(), 10)}
				}
				if isOneway(ctx) {
					if meta == nil {
						meta = map[string]string{}
					}
					meta["one-way"] = "true"
				}
				req.Meta = meta
				req.CalculateHeaderLength()
				req.CalculateBodyLength()
				resp, err := p.Invoke(ctx, req)
				if err != nil {
					return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
				}
				var retErr error
				if len(resp.Error) > 0 {
					retErr = errors.New(string(resp.Error))
				}
				if len(resp.Data) > 0 {
					cpData, er := compress.UnCompression(resp.Data)
					if er != nil {
						return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
					}
					err = serializer.Decode(cpData, retVal.Interface())
					if err != nil {
						return []reflect.Value{reflect.Zero(fieldTyp.Type.Out(0)), reflect.ValueOf(err)}
					}
				}
				if retErr != nil {
					return []reflect.Value{retVal, reflect.ValueOf(retErr)}
				}
				return []reflect.Value{retVal, reflect.Zero(reflect.TypeOf(new(error)).Elem())}
			})
			fieldVal.Set(fn)
		}
	}
	return nil
}

func WithSerializer(sl serialize.Serializer) ClientOpt {
	return func(client *Client) {
		client.serializer = sl
	}
}

func WithCompress(cp compress.Compress) ClientOpt {
	return func(client *Client) {
		client.compress = cp
	}
}

func NewClient(network, add string, opts ...ClientOpt) (*Client, error) {
	p, err := pool.NewChannelPool(&pool.Config{
		InitialCap:  1,
		MaxCap:      30,
		MaxIdle:     10,
		IdleTimeout: time.Minute * 3,
		Factory: func() (interface{}, error) {
			return net.DialTimeout(network, add, time.Minute)
		},
		Close: func(i interface{}) error {
			return i.(net.Conn).Close()
		},
	})
	if err != nil {
		return nil, err
	}
	res := &Client{
		pool:       p,
		serializer: json.Serializer{},
		compress:   gz.GzipCompress{},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

type Client struct {
	pool       pool.Pool
	serializer serialize.Serializer
	compress   compress.Compress
}

func (c *Client) Invoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	ch := make(chan struct{})
	var (
		msg *message.Response
		err error
	)
	go func() {
		msg, err = c.doInvoke(ctx, req)
		close(ch)
	}()
	select {
	case <-ch:
		return msg, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) doInvoke(ctx context.Context, req *message.Request) (*message.Response, error) {
	reqData := message.EncodeReq(req)
	msg, err := c.Send(ctx, reqData)
	if err != nil {
		return nil, err
	}
	return message.DecodeResp(msg), nil
}

func (c *Client) Send(ctx context.Context, data []byte) ([]byte, error) {
	val, err := c.pool.Get()
	if err != nil {
		return nil, err
	}
	defer c.pool.Put(val)
	_, err = val.(net.Conn).Write(data)
	if err != nil {
		return nil, err
	}
	if isOneway(ctx) {
		return nil, errors.New("mirco：这是一个 oneway 调用，你不应该处理结果")
	}
	return ReadMsg(val.(net.Conn))
}
