package cachita

import (
	"errors"
	"fmt"
	"github.com/gadelkareem/go-helpers"
	"reflect"
	"strings"
	"time"
)

type (
	Cache interface {
		Get(key string, i interface{}) error
		Put(key string, i interface{}, ttl time.Duration) error // ttl 0:default ttl, -1: keep forever
		Exists(key string) bool
		Invalidate(key string) error
	}
	cacheError struct {
		err error
	}
	record struct {
		Data      interface{}
		ExpiredAt time.Time
	}
)

var (
	ErrNotFound = newError("cache not found")
	ErrExpired  = newError("cache expired")
)

func Id(params ...string) string {
	return strings.ToLower(helpers.Md5(strings.Join(params, "_")))
}

func newError(msg string) cacheError {
	return cacheError{errors.New(msg)}
}

func (e cacheError) Error() string { return e.err.Error() }

func ExpiredAt(ttl, defaultTtl time.Duration) (expiredAt time.Time) {
	if ttl == 0 {
		expiredAt = time.Now().Add(defaultTtl)
	} else if ttl == -1 {
		expiredAt = time.Now().Add((1000000000) * time.Second) // ten years
	} else {
		expiredAt = time.Now().Add(ttl)
	}
	return
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
	if ok {
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
	v = deReference(v)

	s := reflect.ValueOf(source)
	if !s.IsValid() {
		return errors.New("cachita: source is not valid")
	}
	s = deReference(s)
	v.Set(s)
	return nil
}

func deReference(v reflect.Value) reflect.Value {
	if (v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface) && !v.IsNil() {
		return v.Elem()
	}
	return v
}
