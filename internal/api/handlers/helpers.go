package handlers

import "regexp"

// isValidEmail validates a basic email format.
func isValidEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return reg.MatchString(email)
}
