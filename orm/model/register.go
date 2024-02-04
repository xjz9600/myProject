package orm

import (
	"myProject/orm/internal/errs"
	"myProject/orm/model"
	"reflect"
	"sync"
)

type Register interface {
	Get(entity any) (*model.Model, error)
	Register(entity any, opt ...model.modelOpt) (*model.Model, error)
}

type syncMapRegister struct {
	models sync.Map
}

func NewSyncMapRegister() *syncMapRegister {
	return &syncMapRegister{}
}

type syncRegister struct {
	mutex  sync.RWMutex
	models map[reflect.Type]*model.Model
}

func NewSyncRegister() *syncRegister {
	return &syncRegister{
		models: map[reflect.Type]*model.Model{},
	}
}

func WithTableName(tableName string) model.modelOpt {
	return func(m *model.Model) error {
		m.TableName = tableName
		return nil
	}
}

func WithColumnName(fieldName string, columnName string) model.modelOpt {
	return func(m *model.Model) error {
		if _, ok := m.FieldMap[fieldName]; !ok {
			return errs.NewErrUnknownField(fieldName)
		}
		src := m.FieldMap[fieldName]
		delete(m.ColumnMap, src.ColName)
		src.ColName = columnName
		m.FieldMap[fieldName] = src
		m.ColumnMap[columnName] = src
		return nil
	}
}

func (r *syncMapRegister) Register(entity any, opts ...model.modelOpt) (*model.Model, error) {
	typ := reflect.TypeOf(entity)
	m, err := model.parseModel(entity)
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

func (r *syncMapRegister) Get(entity any) (*model.Model, error) {
	typ := reflect.TypeOf(entity)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model.Model), nil
	}
	m, err := model.parseModel(entity)
	if err != nil {
		return nil, err
	}
	// 会有重复解析的问题但是影响不大
	r.models.Store(typ, m)
	return m.(*model.Model), nil
}

func (r *syncRegister) Get(entity any) (*model.Model, error) {
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
	m, err := model.parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m
	r.mutex.Unlock()
	return m, nil
}

func (r *syncRegister) Register(entity any, opts ...model.modelOpt) (*model.Model, error) {
	typ := reflect.TypeOf(entity)
	m, err := model.parseModel(entity)
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
	m, err = model.parseModel(entity)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m
	r.mutex.Unlock()
	return m, err
}
