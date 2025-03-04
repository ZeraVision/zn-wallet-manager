package functions

import (
	"fmt"
	"log"
)

func SafeString(val interface{}) (string, bool) {
	if str, ok := val.(string); ok {
		return str, true
	}
	log.Printf("[WARN] Expected string but got: %+v\n", val)
	return "", false
}

func SafeInt64(val interface{}) (int64, bool) {
	if num, ok := val.(float64); ok {
		return int64(num), true
	}
	log.Printf("[WARN] Expected int64 but got: %+v\n", val)
	return 0, false
}

func SafeOptionalString(val interface{}) (*string, bool) {
	if str, ok := val.(string); ok {
		return &str, true
	}
	return nil, false
}

func SafeOptionalFloat(val interface{}) (*float64, bool) {
	if num, ok := val.(float64); ok {
		return &num, true
	}
	return nil, false
}

// SafeMapFloatToString extracts a float value from a map and converts it to a string.
// If the key does not exist or is not a float, it returns "0".
func SafeMapFloatToString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if num, ok := val.(float64); ok {
			return fmt.Sprintf("%.8f", num)
		}
	}
	return "0"
}

// SafeMapFloat extracts a float64 value from a map.
// If the key does not exist or is not a float, it returns nil.
func SafeMapFloat(m map[string]interface{}, key string) *float64 {
	if val, exists := m[key]; exists {
		if num, ok := val.(float64); ok {
			return &num
		}
	}
	return nil
}

// SafeMapFloatDefault extracts a float64 value from a map, returning a default if missing.
// If the key does not exist or is not a float, it returns the provided default value.
func SafeMapFloatDefault(m map[string]interface{}, key string, defaultValue float64) float64 {
	if val, exists := m[key]; exists {
		if num, ok := val.(float64); ok {
			return num
		}
	}
	return defaultValue
}

func SafeOptionalBool(val interface{}) (*bool, error) {
	if val == nil {
		return nil, nil
	}

	if b, ok := val.(bool); ok {
		return &b, nil
	}

	return nil, fmt.Errorf("expected bool but got: %+v", val)
}

// SafeMapString extracts a string value from a map safely.
// Returns nil if the key is missing or not a string.
func SafeMapString(m map[string]interface{}, key string) *string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return &str
		}
	}
	return nil
}

// SafeMapStringArray extracts a string slice from a map safely.
// Returns nil if the key is missing or the data type is incorrect.
func SafeMapStringArray(m map[string]interface{}, key string) *[]string {
	if val, exists := m[key]; exists {
		if arr, ok := val.([]interface{}); ok {
			var stringArray []string
			for _, item := range arr {
				if str, valid := item.(string); valid {
					stringArray = append(stringArray, str)
				}
			}
			return &stringArray
		}
	}
	return nil
}

// SafeMapFloatArray extracts a float64 slice from a map safely.
// Returns nil if the key is missing or the data type is incorrect.
func SafeMapFloatArray(m map[string]interface{}, key string) *[]float64 {
	if val, exists := m[key]; exists {
		if arr, ok := val.([]interface{}); ok {
			var floatArray []float64
			for _, item := range arr {
				if num, valid := item.(float64); valid {
					floatArray = append(floatArray, num)
				}
			}
			return &floatArray
		}
	}
	return nil
}

// SafeMapBoolDefault extracts a boolean value from a map safely.
// If the key does not exist or is not a boolean, it returns the provided default value.
func SafeMapBoolDefault(m map[string]interface{}, key string, defaultValue bool) bool {
	if val, exists := m[key]; exists {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}
