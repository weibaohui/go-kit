package datekit

import (
	"testing"
	"time"
)

func TestStringToTime(t *testing.T) {
	parse, err := time.Parse(normalLayout, "2018-01-01 15:15:15")
	if err != nil {
		t.Fatalf(err.Error())
	}
	toTime, err := StringToTime("2018-01-01 15:15:15", normalLayout)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Log(parse.Equal(toTime))
	t.Log(toTime)
}
