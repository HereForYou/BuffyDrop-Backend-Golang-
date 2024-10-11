package utils

import (
	"go-test/models"
	"reflect"
	"regexp"
)

// Function to check if the friends array contains id
func HasFriendWithId(friends []models.Friend, id string) bool {
	for _, friend := range friends {
		if friend.Id == id {
			return true
		}
	}
	return false
}

// normalizeDateString normalizes the fractional seconds in the date string
func NormalizeDateString(dateStr string) (string, error) {
	// Regular expression to match the fractional seconds
	re := regexp.MustCompile(`(\.\d+)(\s+\+\d{4}\s+\w+)`)
	matches := re.FindStringSubmatch(dateStr)

	if len(matches) > 0 {
		// Extract the fractional seconds and timezone part
		fractionalSeconds := matches[1]
		timezone := matches[2]

		// Normalize the fractional seconds to three digits
		switch len(fractionalSeconds) {
		case 2: // e.g., .7
			fractionalSeconds += "00" // Add a zero to make it .740
		case 3: // e.g., .74
			fractionalSeconds += "0" // Add two zeros to make it .700
		}

		// Replace the original fractional seconds in the date string
		normalizedDateStr := re.ReplaceAllString(dateStr, fractionalSeconds+timezone)
		return normalizedDateStr, nil
	}

	return dateStr, nil // Return the original if no fractional seconds are found
}

// Function to check if a slice contains a specific field value or a basic type value
func ContainsValue(slice interface{}, fieldName string, value interface{}) bool {
	v := reflect.ValueOf(slice)

	// Check if the provided slice is indeed a slice
	if v.Kind() != reflect.Slice {
		return false
	}

	// Check if the slice is a slice of structs
	if v.Type().Elem().Kind() == reflect.Struct {
		// Iterate through the slice of structs
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)

			// Get the field by name
			field := item.FieldByName(fieldName)
			if field.IsValid() && field.Interface() == value {
				return true
			}
		}
	} else {
		// Otherwise, assume it's a slice of basic types
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Interface() == value {
				return true
			}
		}
	}
	return false
}
