package sync

import "sync"

type MyBiz struct {
	once sync.Once
}

// Init 使用指针避免负责
func (m *MyBiz) Init() {
	m.once.Do(func() {

	})
}

// MyBizV1 里面的once使用了指针可以直接定义在结构体上
type MyBizV1 struct {
	once *sync.Once
}

func (m MyBizV1) Init() {
	m.once.Do(func() {

	})
}

type MyBusiness interface {
	DoSomething()
}

type singleton struct {
}

func (s *singleton) DoSomething() {
	panic("implement me")
}

var s MyBusiness

var singletonOnce sync.Once

func GetSingleton() MyBusiness {
	singletonOnce.Do(func() {
		s = &singleton{}
	})
	return s
}

func init() {
	// 用包初始化函数取代 once
	s = &singleton{}
}
