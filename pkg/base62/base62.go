package base62

import "strings"

// 62进制转换

// 为了避免被人恶意请求，可以将字符串打乱，如果打乱，则对测试样例进行修改

var (
	base62Str string
)

// MustInit must be called before using base62.
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string")
	}
	base62Str = bs
}

func IntToString(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	var result string
	for seq > 0 {
		mod := seq % 62
		result = string(base62Str[mod]) + result
		seq /= 62
	}
	return result
}

func StringToInt(s string) uint64 {
	var result uint64
	for _, v := range s {
		result = result*62 + uint64(strings.Index(base62Str, string(v)))
	}
	return result
}
