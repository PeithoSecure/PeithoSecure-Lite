package models

type EmailVerificationToken struct {
	ID        int
	Username  string
	Token     string
	ExpiresAt string
}
