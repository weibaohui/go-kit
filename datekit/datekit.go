package datekit

import "time"

func NowString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Year() int {
	y, _, _ := time.Now().Date()
	return y
}
func Month() int {
	_, m, _ := time.Now().Date()
	return int(m)
}
func Day() int {
	_, _, d := time.Now().Date()
	return d
}
