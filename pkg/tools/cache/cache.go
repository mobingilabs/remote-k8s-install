package cache

import (
	"sync"
)

var memoryCache map[string]map[string]interface{}
var mutex sync.RWMutex

func init() {
	memoryCache = make(map[string]map[string]interface{})
}

func GetMap(prefix string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	v, exists := memoryCache[prefix]
	return v, exists
}

func GetOne(prefix, key string) (interface{}, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	v, exists := memoryCache[prefix][key]
	return v, exists
}

func Put(prefix, key string, value interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if _, mapExists := memoryCache[prefix]; !mapExists {
		memoryCache[prefix] = make(map[string]interface{})
	}
	_, exists := memoryCache[prefix][key]
	memoryCache[prefix][key] = value
	return exists
}
