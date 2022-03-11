package queue

import (
	"sync"
	"time"
)

var iSync Sync

func GetSync() Sync {
	if iSync == nil {
		panic("iSync is nil")
	}

	return iSync
}

func ResisterSync(v Sync) {
	if iSync != nil {
		panic("iSync not is nil")
	}
	iSync = v
}

type Sync interface {
	// 添加任务
	AddSync(func() error)
	// 等待所有任务完成
	WaitSync()
}

type SyncImpl struct {
	ch chan func() error
	sync.WaitGroup
}

func NewSync() Sync {
	s := &SyncImpl{}
	s.init()
	s.StartSyncData()

	return s
}

func (s *SyncImpl) AddSync(f func() error) {

	// 增加次数
	s.Add(1)
	// 完成后扣减次数
	s.ch <- func() error {
		defer s.Done()
		for {
			err := f()
			if err != nil {
				log.Error("sync errmsg: %+v", err)
				time.Sleep(time.Second / 2)
			} else {
				return nil
			}
		}
	}
}

// 初始数据结构
func (s *SyncImpl) init() {
	s.ch = make(chan func() error, 1024)
}

// 启动同步数据
func (s *SyncImpl) StartSyncData() {
	go s._SyncData()
}

// 同步数据
func (s *SyncImpl) _SyncData() {
	for f := range s.ch {
		f()
	}
}

func (s *SyncImpl) WaitSync() {

	// 等待任务完成
	s.Wait()
}
