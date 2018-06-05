package rediskit

import (
	"encoding/json"
	"strings"
	"time"
)

// 检查是否存在
func (r *redisKit) IsExists(key string) bool {
	result, err := r.DB.Exists(key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

// 获取字符串，没有的话先用dataLoader加载数据
func (r *redisKit) GetStringWithDataLoader(key string, dataLoader func() (string, error), duration time.Duration) (string, error) {

	if r.IsExists(key) {
		return r.GetString(key)
	} else {

		ret, err := dataLoader()
		if err != nil {
			return "", err
		}
		r.Set(key, ret, duration)
		return ret, err
	}
}

// 获取字符串
func (r *redisKit) GetString(key string) (string, error) {
	b, err := r.DB.Get(key).Result()
	if err != nil {
		return "", err
	}
	b = strings.Replace(b, "\"", "", -1)
	return b, nil
}

func (r *redisKit) Delete(key string) error {
	_, err := r.DB.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}
func (r *redisKit) Get(key string, clazz interface{}) error {

	b, err := r.DB.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(b), clazz)
	return err
}

func (r *redisKit) Set(key string, clazz interface{}, duration time.Duration) error {

	b, err := json.Marshal(clazz)
	if err != nil {
		return err
	}
	return r.DB.Set(key, string(b), duration).Err()

}

func (r *redisKit) SetOneHour(key string, clazz interface{}) error {
	return r.Set(key, clazz, time.Duration(1)*time.Hour)
}

func (r *redisKit) SetOneDay(key string, clazz interface{}) error {
	return r.Set(key, clazz, time.Duration(1)*time.Hour*24)
}

func (r *redisKit) SetTenMinute(key string, clazz interface{}) error {
	return r.Set(key, clazz, time.Duration(1)*time.Minute*10)
}

func (r *redisKit) SetWithDuration(key string, clazz interface{}, duration time.Duration) error {
	return r.Set(key, clazz, duration)
}
