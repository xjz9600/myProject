package orm

import (
	"myProject/orm/internal/errs"
	"reflect"
	"sync"
)

type Register interface {
	Get(entity any) (*model, error)
	Register(entity any, opt ...modelOpt) (*model, error)
}

type syncMapRegister struct {
	models sync.Map
}

func NewSyncMapRegister() *syncMapRegister {
	return &syncMapRegister{}
}

type syncRegister struct {
	mutex  sync.RWMutex
	models map[reflect.Type]*model
}

func NewSyncRegister() *syncRegister {
	return &syncRegister{
		models: map[reflect.Type]*model{},
	}
}

func WithTableName(tableName string) modelOpt {
	return func(m *model) error {
		m.tableName = tableName
		return nil
	}
}

func WithColumnName(fieldName string, columnName string) modelOpt {
	return func(m *model) error {
		if _, ok := m.fieldMap[fieldName]; !ok {
			return errs.NewErrUnknownField(fieldName)
		}
		m.fieldMap[fieldName].colName = columnName
		return nil
	}
}

func (r *syncMapRegister) Register(entity any, opts ...modelOpt) (*model, error) {
	typ := reflect.TypeOf(entity)
	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	r.models.Store(typ, m)
	return m, err
}

func (r *syncMapRegister) Get(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}
	// 会有重复解析的问题但是影响不大
	r.models.Store(typ, m)
	return m.(*model), nil
}

func (r *syncRegister) Get(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	r.mutex.RLock()
	m, ok := r.models[typ]
	if ok {
		return m, nil
	}
	r.mutex.RUnlock()
	// double-check
	r.mutex.Lock()
	m, ok = r.models[typ]
	if ok {
		return m, nil
	}
	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m
	r.mutex.Unlock()
	return m, nil
}

func (r *syncRegister) Register(entity any, opts ...modelOpt) (*model, error) {
	typ := reflect.TypeOf(entity)
	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	r.mutex.Lock()
	m, ok := r.models[typ]
	if ok {
		return m, nil
	}
	m, err = parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m
	r.mutex.Unlock()
	return m, err
}
