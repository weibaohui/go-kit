package propkit

import (
	"testing"
)

func TestConfig_Get(t *testing.T) {
	x := Init().Use("./src/config/config.json").Get("server.port")
	t.Log(x)
}
