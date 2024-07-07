package convert

import (
	"strconv"
	"unsafe"
)

func IntToStr[T ~int | ~int32 | ~int64](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

/*
byte to string with zero copy

@warning
use unsafe
전달되는 b의 데이터가 수정이 되는 경우 좋지 않음.
*/
func ZeroCopyByteToString(b []byte) string {
	return *((*string)(unsafe.Pointer(&b)))
}

/*
string to bytes

@waring
use unsafe
*/
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

/*
Map 으로 변환하기

@use

	convert.ToMap([]interface{}, func(data interface{}) string {
		return <key>
	})
*/
func ToMap[T any, V comparable](s []T, f func(T) V) map[V]T {
	ret := make(map[V]T, len(s))
	for _, e := range s {
		ret[f(e)] = e
	}
	return ret
}

/*
Map 으로 변환하기. Value 가 List

@use

	convert.ToMap([]interface{}, func(data interface{}) string {
		return <key>
	})
*/
func ToMapSli[T any, V comparable](s []T, f func(T) V) map[V][]T {
	ret := make(map[V][]T, len(s))
	for _, e := range s {
		if ret[f(e)] == nil {
			ret[f(e)] = make([]T, 0, 10)
		}
		ret[f(e)] = append(ret[f(e)], e)
	}
	return ret
}

/*
Slice 로 변환하기

@use
*/
func ToSli[T any, V comparable](m map[V]T) []T {
	ret := make([]T, 0, len(m))
	for _, e := range m {

		ret = append(ret, e)
	}
	return ret
}
