package utils

import (
	"encoding/json"
	"sync"
)

// 并发安全容器工具
type Map struct {
	sync.RWMutex
	m map[interface{}]interface{}
}

func (m *Map) init() { // 延迟初始化。只有在第一次写入时才创建底层的 map，节省内存
	if m.m == nil {
		m.m = make(map[interface{}]interface{})
	}
}

func (m *Map) UnsafeGet(key interface{}) interface{} {
	if m.m == nil {
		return nil
	} else {
		return m.m[key]
	}
}

func (m *Map) Get(key interface{}) interface{} {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeGet(key)
}

func (m *Map) UnsafeSet(key interface{}, value interface{}) {
	m.init()
	m.m[key] = value
}

func (m *Map) Set(key interface{}, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeSet(key, value)
}

// 这是非常实用的功能。它检查某个 Key 是否存在，如果不存在则设置新值并返回 nil；如果已存在则返回旧值。
// 在实现类似“单机版分布式锁”或防止重复初始化时非常有用。
func (m *Map) TestAndSet(key interface{}, value interface{}) interface{} {
	m.Lock()
	defer m.Unlock()

	m.init()

	if v, ok := m.m[key]; ok {
		return v
	} else {
		m.m[key] = value
		return nil
	}
}

func (m *Map) UnsafeDel(key interface{}) {
	m.init()
	delete(m.m, key)
}

func (m *Map) Del(key interface{}) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeDel(key)
}

func (m *Map) UnsafeLen() int {
	if m.m == nil {
		return 0
	} else {
		return len(m.m)
	}
}

func (m *Map) Len() int {
	m.RLock()
	defer m.RUnlock()
	return m.UnsafeLen()
}

func (m *Map) UnsafeRange(f func(interface{}, interface{})) {
	if m.m == nil {
		return
	}
	for k, v := range m.m {
		f(k, v)
	}
}

// 遍历时只加读锁，允许其他协程同时读，但不允许写
func (m *Map) RLockRange(f func(interface{}, interface{})) {
	m.RLock()
	defer m.RUnlock()
	m.UnsafeRange(f)
}

// 遍历时加写锁，完全独占，防止遍历过程中数据被修改
func (m *Map) LockRange(f func(interface{}, interface{})) {
	m.Lock()
	defer m.Unlock()
	m.UnsafeRange(f)
}

// 将 map[string]interface{} 转换为 JSON 字符串
func MapToJsonString(param map[string]interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

// 将 JSON 字符串解析回 Map
func JsonStringToMap(str string) (tempMap map[string]interface{}) {
	_ = json.Unmarshal([]byte(str), &tempMap)
	return tempMap
}

// 一个业务逻辑函数。看起来是用于从配置项中读取开关状态（默认为 true，除非显式设为非 1 的值）
func GetSwitchFromOptions(Options map[string]interface{}, key string) (result bool) {
	if flag, ok := Options[key]; !ok || flag == 1 {
		return true
	}
	return false
}
