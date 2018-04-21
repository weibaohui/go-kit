package filekit

import "testing"

func TestCurrentPath(t *testing.T) {
	s := CurrentPath()
	t.Log(s)
}
