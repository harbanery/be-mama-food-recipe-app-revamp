package helper

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type WebResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

func EmptyObject() interface{} {
	return make(map[string]interface{})
}

func EmptyArray() interface{} {
	return []interface{}{}
}

func UsernameFromEmail(email string) string {
	return strings.Split(email, "@")[0]
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

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
