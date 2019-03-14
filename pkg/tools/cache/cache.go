package cache

import (
	"sync"
)

var memoryCache map[string]interface{}
var mutex sync.RWMutex

func init() {
	memoryCache = make(map[string]interface{})
}

func Get(key string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	v, exists := memoryCache[key]
	return v, exists
}

func Put(key string, value interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := memoryCache[key]
	memoryCache[key] = value
	return exists
}
