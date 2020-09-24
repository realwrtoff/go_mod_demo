package cache

import "sync"

type KvStore interface {
	Get(key interface{}) (interface{}, bool)
	Set(key interface{}, value interface{}) bool
}

type MemKv struct {
	kvs   map[interface{}]interface{}
	mutex sync.RWMutex
}

type MgoKv struct {
	kvs *Mongo
}