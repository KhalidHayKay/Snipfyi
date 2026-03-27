package utils

import "math/rand"

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Encode(num int64) string {
	if num == 0 {
		return "0"
	}

	base := int64(len(charset))
	result := ""

	for num > 0 {
		rem := num % base
		result = string(charset[rem]) + result
		num = num / base
	}

	return result
}

func EncodeWithPadding(id int64) string {
	base62 := Encode(id) // your existing Base62 encoder
	if len(base62) >= 2 {
		return base62
	}

	for len(base62) < 3 {
		randIndex := rand.Intn(len(charset))
		base62 = base62 + string(charset[randIndex])
	}

	return base62
}
