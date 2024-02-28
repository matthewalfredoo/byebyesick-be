package util

import (
	"crypto/rand"
	"fmt"
	"halodeksik-be/app/appconstant"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func IsEmptyString(str string) bool {
	return str == ""
}

func GetCurrentDateAndTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)
}

func GetCurrentDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func ParseDateTime(timeStr string, timeFormat ...string) (time.Time, error) {
	if len(timeFormat) > 0 {
		return time.Parse(timeFormat[0], timeStr)
	}
	return time.Parse(appconstant.TimeFormatQueryParam, timeStr)
}

func RandomToken(marker string) (string, error) {
	bTime, _ := time.Now().MarshalText()
	bMarker := []byte(marker)

	b := append(bMarker, bTime...)
	_, err := rand.Read(b)

	return fmt.Sprintf("%x", b), err
}

func ParseInt64(str string) (int64, error) {
	return strconv.ParseInt(str, 10, 64)
}

func PascalToSnake(input string) string {
	var builder strings.Builder

	for i, char := range input {
		if unicode.IsUpper(char) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(char))
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func AppendAtIndex[T any](arr []T, index int, value T) []T {
	if len(arr) == 0 || arr == nil {
		arr = make([]T, 0)
		arr = append(arr, value)
		return arr
	}

	if index > len(arr)-1 {
		index = len(arr) - 1
	}

	newArr := make([]T, 0)
	newArr = append(newArr, arr[:index]...)
	newArr = append(newArr, value)
	newArr = append(newArr, arr[index:]...)

	return newArr
}
