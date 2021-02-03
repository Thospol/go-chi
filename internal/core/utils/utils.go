package utils

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	numberOfShortPostText = 100
)

// GenerateUniqueCode gengerate unique code
func GenerateUniqueCode() string {
	return time.Now().In(LoadLocation()).Format("20060102150405")
}

// GenerateSeqNo generate seq no
func GenerateSeqNo(n, length int) string {
	digits := CountDigits(n)
	if digits > length {
		return fmt.Sprintf("%d", n)
	}

	return fmt.Sprintf("%0*d", length, n)
}

// CountDigits count digits
func CountDigits(i int) (count int) {
	for i != 0 {
		i /= 10
		count = count + 1
	}
	return count
}

// GenerateNumber generate number
func GenerateNumber(number, length int) string {
	if number <= 0 {
		return ""
	}

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", length), number)
}

// UniqueSliceString for remove duplicate from slice string
func UniqueSliceString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if entry != "" {
				list = append(list, entry)
			}
		}
	}

	return list
}

// UniqueSliceInt for remove duplicate from slice int
func UniqueSliceInt(slice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if entry != 0 {
				list = append(list, entry)
			}
		}
	}

	return list
}

// UniqueSliceUInt for remove duplicate from slice uint
func UniqueSliceUInt(slice []uint) []uint {
	keys := make(map[uint]bool)
	list := []uint{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if entry != 0 {
				list = append(list, entry)
			}
		}
	}

	return list
}

// IsSliceStringChanged is slice string change
func IsSliceStringChanged(original, compare []string) bool {
	if len(original) != len(compare) {
		return true
	}

	var count int
	for _, r := range compare {
		var i int
		for i <= len(original)-1 {
			if r == original[i] {
				count++
				break
			}
			i++
		}
	}

	return count != len(original)
}

// IsEmptySlice is empty slice
func IsEmptySlice(values interface{}) bool {
	if values, ok := values.([]int); ok {
		for _, value := range values {
			if value == 0 {
				return true
			}
		}
		return false
	}

	if values, ok := values.([]string); ok {
		for _, value := range values {
			if value == "" {
				return true
			}
		}
		return false
	}

	return true
}

// RoundUp is round up decimal position at 1 if value more than 0
func RoundUp(v float64) float64 {
	return float64(int((v*100)+9)/10) / 10.0
}

// WrapPassword for wrap password Example: ********
func WrapPassword(password string) string {
	var wrapPassword string
	lengthPassword := utf8.RuneCountInString(password)
	for i := 0; i < lengthPassword; i++ {
		wrapPassword = wrapPassword + "*"
	}

	return wrapPassword
}

// GetInitialName for get first character.
func GetInitialName(name string) string {
	runeName := []rune(name)
	if IsValidEnglishAlphabet(string(runeName[0])) {
		name = strings.ToUpper(string(runeName[0]))
		return name
	}

	return string(runeName[0])
}

// IsValidEmails is valid emails
func IsValidEmails(emails []string) bool {
	for _, email := range emails {
		if email == "" || !IsValidEmail(email) {
			return false
		}
	}

	return true
}

// TrimSpaces trim slice
func TrimSpaces(slice []string) []string {
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}

	return slice
}

// ToLower to lower
func ToLower(slice []string) []string {
	for i := range slice {
		slice[i] = strings.ToLower(slice[i])
	}

	return slice
}

// RoundFloat64 round float64 (100 = .2f)
func RoundFloat64(x, unit float64) float64 {
	return math.Round(x*unit) / unit
}

// SubString sub string
func SubString(sourceText string, showDot bool) string {
	newString := sourceText
	if len([]rune(sourceText)) > numberOfShortPostText {
		dotString := ""
		if showDot {
			dotString = "..."
		}
		newString = fmt.Sprintf("%s%s", string([]rune(sourceText)[:numberOfShortPostText]), dotString)
	}

	return newString
}

// ConvertToArrayInterface convert to array interface...
func ConvertToArrayInterface(x interface{}) []interface{} {
	result := []interface{}{}
	switch reflect.TypeOf(x).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(x)
		for i := 0; i < s.Len(); i++ {
			result = append(result, s.Index(i).Interface())
		}
	}

	return result
}
