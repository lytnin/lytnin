package services

import (
	"ppmeweb/lytnin"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// KVStore service provides key/value store access to the application
type KVStore struct {
	Conn redis.Conn
}

// Info returns information about the key/value store
func (s *KVStore) Info() interface{} {
	// kv store
	tmp, _ := s.Conn.Do("INFO")
	str := string(tmp.([]byte))
	info := strings.Split(str, "\r\n")
	return info[1]
}

// Init initializes the key/value service and registers it with the application
func (s *KVStore) Init(a *lytnin.Application) {
	r, err := redis.DialURL(a.Config.Get("REDIS_URL"))
	checkErr(err)
	s.Conn = r

	a.AddService("kvstore", s)
}

// Close releases any resources used by the service
func (s *KVStore) Close() {
	s.Conn.Close()
}
