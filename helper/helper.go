package helper

import (
	"strings"

	"github.com/google/uuid"
)

type WebResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func EmptyObject() interface{} {
	return make(map[string]interface{})
}

func EmptyArray() interface{} {
	return []interface{}{}
}

func GenerateSessionID() string {
	randomUUID := uuid.New()
	return strings.Replace(randomUUID.String(), "-", "", -1)
}

func ToUppercase(s *string) *string {
	if s == nil {
		return nil
	}
	upper := strings.ToUpper(*s)
	return &upper
}
