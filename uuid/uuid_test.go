package uuid

import (
	"fmt"
	"testing"
)

func TestUUID(t *testing.T) {
	m := make(map[string]int)
	for i := 0; i < 10000; i++ {

		uuid := UUID32()
		fmt.Printf("%4d = %s \n", i, uuid)
		if _, ok := m[uuid]; ok {
			t.Fatal("发现重复" + uuid)
		}
		m[uuid] = 1
	}
}
