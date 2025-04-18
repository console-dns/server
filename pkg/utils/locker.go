package utils

import "sync"

// DataRwLocker 线程安全的数据管理
type DataRwLocker[T any] struct {
	data T            // 数据
	lock sync.RWMutex // 读写锁
}

func NewDataRwLocker[T any](data T) *DataRwLocker[T] {
	return &DataRwLocker[T]{
		data: data,
		lock: sync.RWMutex{},
	}
}

func (d *DataRwLocker[T]) ReadOnly(roFunc func(T)) {
	defer func() {
		d.lock.RUnlock()
	}()
	d.lock.RLock()
	roFunc(d.data)
}

func (d *DataRwLocker[T]) ReadWrite(roFunc func(T)) {
	defer func() {
		d.lock.Unlock()
	}()
	d.lock.Lock()
	roFunc(d.data)
}

// WithReadOnly 以只读模式打开数据 （注意及时关闭锁）
func (d *DataRwLocker[T]) WithReadOnly() (result T, unlock func()) {
	unlock = func() {
		d.lock.RUnlock()
	}
	d.lock.RLock()
	result = d.data
	return
}

// WithReadWrite 以写入的方式打开锁 （注意及时关闭锁）
func (d *DataRwLocker[T]) WithReadWrite() (result T, unlock func()) {
	unlock = func() {
		d.lock.Unlock()
	}
	d.lock.Lock()
	result = d.data
	return
}
