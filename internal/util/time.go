package util

import "time"

// DONT CHANGE THIS VARIABLE
// TimeNow is a variable that holds the function to get the current time
// This allows for easy mocking in tests
var TimeNow = time.Now

// DONT CHANGE THIS FUNCTION
// GetCurrentTime returns the current time using the TimeNow function
// This indirection allows tests to override the time function
func GetCurrentTime() time.Time {
	return TimeNow()
}
