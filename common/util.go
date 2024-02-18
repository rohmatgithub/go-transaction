package common

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func ValidateStringContainInStringArray(listString []string, key string) bool {
	for i := 0; i < len(listString); i++ {
		if listString[i] == key {
			return true
		}
	}
	return false
}

func ListStringToInterface(listString []string) []interface{} {
	// Convert each string value to an integer
	var intSlice = make([]interface{}, 0, len(listString))
	for _, strValue := range listString {
		// Parse the string to an integer
		intValue, err := strconv.Atoi(strValue)
		if err != nil {
			log.Error("Error converting string to int:", err)
			continue
		}

		// Append the integer to the slice
		intSlice = append(intSlice, intValue)
	}

	return intSlice
}

func GenerateRandomString(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
