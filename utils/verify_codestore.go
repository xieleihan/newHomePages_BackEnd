package utils

import (
	"gin/db"
	"time"
)

type RedisStore struct{}

func (r RedisStore) Set(id, value string) error {
	key := "captcha:" + id
	return db.RDB.Set(db.Ctx, key, value, 3*time.Minute).Err()
}

func (r RedisStore) Get(id string, clear bool) string {
	key := "captcha:" + id
	val, err := db.RDB.Get(db.Ctx, key).Result()
	if err != nil {
		return ""
	}
	if clear {
		db.RDB.Del(db.Ctx, key)
	}
	return val
}

func (r RedisStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, clear)
	return val == answer
}
