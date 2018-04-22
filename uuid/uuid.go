package uuid

import (
	"github.com/google/uuid"
	"strings"
)

func UUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return id.String()
}

func UUID32() string {
	id, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return strings.Replace(id.String(), "-", "", -1)
}
