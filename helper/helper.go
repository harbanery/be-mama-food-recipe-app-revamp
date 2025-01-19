package helper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WebResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type RecordCount struct {
	FilteredData int64
	TotalData    int64
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

func ToSlug(title, author string) string {
	// Normalize title and author
	normalize := func(input string) string {
		// Convert to lowercase
		result := strings.ToLower(input)
		// Replace non-alphanumeric characters with a hyphen
		re := regexp.MustCompile(`[^a-z0-9]+`)
		result = re.ReplaceAllString(result, "-")
		// Trim any leading or trailing hyphens
		result = strings.Trim(result, "-")
		return result
	}

	// Normalize title and author
	slugTitle := normalize(title)
	slugAuthor := normalize(author)

	// Combine into the desired format
	return fmt.Sprintf("%s-by-%s", slugTitle, slugAuthor)
}
