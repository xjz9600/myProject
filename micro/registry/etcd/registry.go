package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"myProject/micro/registry"
	"sync"
)

//type Registry struct {
//	cc     *clientv3.Client
//	sess   *concurrency.Session
//	cancel []func()
//	mutex  sync.Mutex
//}
//
//func NewRegistry(cc *clientv3.Client) (*Registry, error) {
//	// 默认过期时间为60s
//	session, err := concurrency.NewSession(cc)
//	if err != nil {
//		return nil, err
//	}
//	return &Registry{
//		cc:   cc,
//		sess: session,
//	}, nil
//}
//
//func (r *Registry) Register(ctx context.Context, si registry.ServiceInstance) error {
//	val, err := json.Marshal(si)
//	if err != nil {
//		return err
//	}
//	_, err = r.cc.Put(ctx, r.instanceKey(si), string(val), clientv3.WithLease(r.sess.Lease()))
//	return err
//}
//
//func (r *Registry) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
//	_, err := r.cc.Delete(ctx, r.instanceKey(si), clientv3.WithLease(r.sess.Lease()))
//	return err
//}
//
//func (r *Registry) ListServices(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
//	resp, err := r.cc.Get(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
//	if err != nil {
//		return nil, err
//	}
//	instances := make([]registry.ServiceInstance, 0, len(resp.Kvs))
//	for _, kv := range resp.Kvs {
//		instance := registry.ServiceInstance{}
//		err = json.Unmarshal(kv.Value, &instance)
//		if err != nil {
//			return nil, err
//		}
//		instances = append(instances, instance)
//	}
//	return instances, nil
//}
//
//func (r *Registry) Subscribe(serviceName string) (<-chan registry.Event, error) {
//	ctx, cancel := context.WithCancel(context.Background())
//	ctx = clientv3.WithRequireLeader(ctx)
//	r.mutex.Lock()
//	r.cancel = append(r.cancel, cancel)
//	r.mutex.Unlock()
//	watchChan := r.cc.Watch(ctx, r.serviceKey(serviceName), clientv3.WithPrefix())
//	res := make(chan registry.Event)
//	go func() {
//		for {
//			select {
//			case resp := <-watchChan:
//				if resp.Err() != nil {
//					continue
//				}
//				if resp.Canceled {
//					return
//				}
//				res <- registry.Event{}
//			case <-ctx.Done():
//				return
//			}
//		}
//	}()
//	return res, nil
//}
//
//func (r *Registry) Close() error {
//	r.mutex.Lock()
//	cancels := r.cancel
//	r.cancel = nil
//	r.mutex.Unlock()
//	for _, c := range cancels {
//		c()
//	}
//	err := r.sess.Close()
//	return err
//}
//
//func (r *Registry) instanceKey(si registry.ServiceInstance) string {
//	return fmt.Sprintf("/micro/%s/%s", si.Name, si.Addr)
//}
//
//func (r *Registry) serviceKey(sn string) string {
//	return fmt.Sprintf("/micro/%s", sn)
//}

type Registry struct {
	c       *clientv3.Client
	sess    *concurrency.Session
	mutex   sync.Mutex
	cancels []func()
}

func NewRegistry(c *clientv3.Client) (*Registry, error) {
	sess, err := concurrency.NewSession(c)
	if err != nil {
		return nil, err
	}
	return &Registry{
		c:    c,
		sess: sess,
	}, nil
}

func (r *Registry) Register(ctx context.Context, si registry.ServiceInstance) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}
	_, err = r.c.Put(ctx, r.instanceKey(si), string(val), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.c.Delete(ctx, r.instanceKey(si), clientv3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) ListServices(ctx context.Context, name string) ([]registry.ServiceInstance, error) {
	getResp, err := r.c.Get(ctx, r.serviceKey(name), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	var res []registry.ServiceInstance
	for _, kv := range getResp.Kvs {
		var si registry.ServiceInstance
		err = json.Unmarshal(kv.Value, &si)
		if err != nil {
			return nil, err
		}
		res = append(res, si)
	}
	return res, nil
}

func (r *Registry) Subscribe(name string) (<-chan registry.Event, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = clientv3.WithRequireLeader(ctx)
	ch := r.c.Watch(ctx, r.serviceKey(name), clientv3.WithPrefix())
	res := make(chan registry.Event)
	r.mutex.Lock()
	r.cancels = append(r.cancels, cancel)
	r.mutex.Unlock()
	go func() {
		for {
			select {
			case resp := <-ch:
				if resp.Err() != nil {
					continue
				}
				if resp.Canceled {
					return
				}
				res <- registry.Event{}
			case <-ctx.Done():
				return
			}
		}
	}()
	return res, nil
}

func (r *Registry) Close() error {
	r.mutex.Lock()
	cancels := r.cancels
	r.cancels = nil
	r.mutex.Unlock()
	for _, ca := range cancels {
		ca()
	}
	return r.sess.Close()
}

func (r *Registry) instanceKey(si registry.ServiceInstance) string {
	return fmt.Sprintf("/micro/%s/%s", si.Name, si.Addr)
}
func (r *Registry) serviceKey(serverName string) string {
	return fmt.Sprintf("/micro/%s", serverName)
}
