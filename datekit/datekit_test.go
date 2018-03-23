package datekit

import (
	"testing"
	"time"
)

func today() int {
	return time.Now().Day()
}
func TestDay(t *testing.T) {
	tests := []struct {
		day int
		ans int
	}{
		{23, 23},
	}
	for _, v := range tests {
		acu := Day()
		if acu != v.day {
			t.Fatalf("day test faild,got %d,except %d", acu, v.ans)
		}
	}
}

func BenchmarkDay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		day := Day()
		if today() != day {
			b.Fatalf("day test faild,got %d,except %d", day, today())
		}
	}
}
