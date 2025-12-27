package config

import "strconv"

func sliceToSetInt(values []int) map[int]struct{} {
	m := make(map[int]struct{}, len(values))
	for _, v := range values {
		m[v] = struct{}{}
	}
	return m
}

func atoiSafe(s string) int {
	v, _ := strconv.ParseInt(s, 10, 64)
	return int(v)
}

func atofSafe(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return f
}
