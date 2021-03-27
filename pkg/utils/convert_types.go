package utils

// String returns a pointer to the string value passed in
func String(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or "" if the pointer is nil
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// String returns a pointer to the string value passed in or nil if the value is empty
func EmptyToNullString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// Int returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// IntValue returns the value of the int pointer passed in or 0 if the pointer is nil.
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int64 returns a pointer to the int64 value passed in
func Int64(v int64) *int64 {
	return &v
}

// Int64Value returns the value of the int64 pointer passed in or 0 if the pointer is nil
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}
