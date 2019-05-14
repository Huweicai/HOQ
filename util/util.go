package ut

import (
	"math/rand"
	"reflect"
	"time"
)

/**
求公共前缀
返回两个字符串公共的序号
0 代表无重复部分
返回第一个不相同字符的序号
*/
func CommonPrefix(a, b string) int {
	s := Min(len(a), len(b))
	for i := 0; i < s; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return s
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

//采样的do something
//rate : 采样率 [0，1]
func DownSamplingDo(rate float32, do func()) {
	if rand.Float32() < rate {
		do()
	}
}

/**
is targets contains target
Contain("a" , ["a","b"]) = true
support both array and map
keywords：包含 集合
*/
func Contain(target interface{}, targets interface{}) bool {
	if target == nil || targets == nil {
		return false
	}
	targetValue := reflect.ValueOf(targets)
	switch reflect.TypeOf(targets).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == target {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(target)).IsValid() {
			return true
		}
	}

	return false
}

/**
Return the intersection of two array
keywords：集合求交集 公共部分
*/
func Common(targetsA []string, targetsB []string) (common []string) {
	for _, a := range targetsA {
		for _, b := range targetsB {
			if a == b {
				common = append(common, a)
				break
			}
		}
	}
	return
}

/**
数组a -b
*/
func SliceReduce(a []string, b []string) (dslice []string) {
	for _, s := range a {
		if !Contain(s, b) {
			dslice = append(dslice, s)
		}
	}
	return
}

/**
a simple util wrap for time.Now()
cause function can not be addressed
*/
func Now() *time.Time {
	t := time.Now()
	return &t
}

/**
for unused error
but doing nothing
*/
func Nothing(i ...interface{}) {

}
