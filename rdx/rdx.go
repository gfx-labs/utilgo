package rdx

import (
	"context"

	"gfx.cafe/open/gokv"
	"gfx.cafe/open/gokv/encoding"
	"gfx.cafe/open/gokv/util"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v9"
)

const kv_name = "gokv"

var mr *miniredis.Miniredis
var mrstore *Store

var _ (gokv.Store) = (*Store)(nil)

type Store struct {
	c     redis.UniversalClient
	codec encoding.Codec
}

func Mem() *Store {
	return NewMiniStore()
}

func NewMiniStore() *Store {
	if mr == nil {
		mr = miniredis.NewMiniRedis()
		mr.Start()
		mrstore = &Store{
			codec: encoding.JSON,
			c: redis.NewClient(&redis.Options{
				Network:  "tcp",
				Addr:     mr.Addr(),
				Username: "",
				Password: "",
				Limiter:  nil,
			})}
	}
	return mrstore
}

func (s *Store) C() redis.UniversalClient {
	return s.c
}

func (s *Store) Set(k string, v interface{}) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}

	data, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	err = s.c.HSet(context.Background(), kv_name, k, data).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(k string, v interface{}) (found bool, err error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	dataString, err := s.c.HGet(context.Background(), kv_name, k).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return true, s.codec.Unmarshal([]byte(dataString), v)
}

func (s *Store) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}

	_, err := s.c.HDel(context.Background(), kv_name, k).Result()
	return err
}

func (s *Store) Close() error {
	return nil
}
