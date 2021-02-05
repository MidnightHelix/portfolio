package base62

import (
	"math"
	"strings"
)

const base = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const b = 62

func Encode(number int) string {
	r := number % b
	res := string(base[r])
	div := number / b
	q := int(math.Floor(float64(div)))

	for q != 0 {
		r = q % b
		temp := q / b
		q = int(math.Floor(float64(temp)))
		res = string(base[int(r)]) + res
	}
	return string(res)
}

func Decode(str string) int {
	res := 0
	for _, r := range str {
		res = (b * res) + strings.Index(base, string(r))
	}
	return res
}
