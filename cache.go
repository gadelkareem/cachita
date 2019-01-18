package cachita

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type (
	Cache interface {
		Get(key string, i interface{}) error
		Put(key string, i interface{}, ttl time.Duration) error // ttl 0:default ttl, -1: keep forever
		Incr(key string, ttl time.Duration) (int64, error)
		Exists(key string) bool
		Invalidate(key string) error
	}
	record struct {
		Data      interface{}
		ExpiredAt time.Time
	}
)

var (
	ErrNotFound = errors.New("cachita: cache not found")
	ErrExpired  = errors.New("cachita: cache expired")
)

func calculateTtl(ttl, defaultTtl time.Duration) time.Duration {
	if ttl == 0 {
		return defaultTtl
	} else if ttl == -1 {
		return 100000 * time.Hour // 11 years
	}
	return ttl
}

func expiredAt(ttl, defaultTtl time.Duration) time.Time {
	return time.Now().Add(calculateTtl(ttl, defaultTtl))
}

func IsErrorOk(err error) bool {
	return err == ErrNotFound || err == ErrExpired
}

func Id(params ...string) string {
	s := strings.Join(params, "_")
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}

func runEvery(ttl time.Duration, f func()) {
	ticker := time.NewTicker(ttl)
	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			}
		}
	}()
}

func TypeAssert(source, target interface{}) (err error) {
	if source == nil {
		return nil
	}
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("cachita: %v", r)
			}
		}
	}()

	if directTypeAssert(source, target) {
		return nil
	}

	v := reflect.ValueOf(target)
	if !v.IsValid() {
		return errors.New("cachita: target is nil")
	}
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("cachita: target is not a settable %T", target)
	}
	v = v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("cachita: target is not a settable %T", target)
	}

	s := reflect.ValueOf(source)
	if !s.IsValid() {
		return errors.New("cachita: source is not valid")
	}
	s = deReference(s)
	v.Set(s)
	return nil
}

func directTypeAssert(source, target interface{}) bool {
	var ok bool
	switch v := target.(type) {
	case *string:
		*v, ok = source.(string)
	case *[]byte:
		*v, ok = source.([]byte)
	case *int:
		*v, ok = source.(int)
	case *int8:
		*v, ok = source.(int8)
	case *int16:
		*v, ok = source.(int16)
	case *int32:
		*v, ok = source.(int32)
	case *int64:
		*v, ok = source.(int64)
	case *uint:
		*v, ok = source.(uint)
	case *uint8:
		*v, ok = source.(uint8)
	case *uint16:
		*v, ok = source.(uint16)
	case *uint32:
		*v, ok = source.(uint32)
	case *uint64:
		*v, ok = source.(uint64)
	case *bool:
		*v, ok = source.(bool)
	case *float32:
		*v, ok = source.(float32)
	case *float64:
		*v, ok = source.(float64)
	case *time.Duration:
		*v, ok = source.(time.Duration)
	case *time.Time:
		*v, ok = source.(time.Time)
	case *[]string:
		*v, ok = source.([]string)
	case *map[string]string:
		*v, ok = source.(map[string]string)
	case *map[string]interface{}:
		*v, ok = source.(map[string]interface{})
	}
	return ok
}

func deReference(v reflect.Value) reflect.Value {
	if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		return v.Elem()
	}
	return v
}
