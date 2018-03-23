package datekit

import "time"

const (
	normalLayout = "2006-01-02 15:04:05"
)

func NowString() string {
	return time.Now().Format(normalLayout)
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
