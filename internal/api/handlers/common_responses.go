package handlers

// UnlockSuccessResponse defines a response for successful unlock
// @Description Response for successful license unlock
type UnlockSuccessResponse struct {
	Message string `json:"message" example:"unlocked"`
}

// GenericMessageResponse is a reusable generic success message
// @Description Standard success message response
type GenericMessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// GenericErrorResponse represents a detailed structured error payload
// @Description Trace error structure with audit details
type GenericErrorResponse struct {
	Message   string `json:"message" example:"Invalid request"`
	Code      string `json:"code" example:"BadRequest"`
	Event     string `json:"event" example:"email_invalid"`
	Actor     string `json:"actor" example:"USER"`
	Timestamp string `json:"timestamp" example:"2025-05-16T10:00:00Z"`
	Debug     string `json:"debug,omitempty" example:"stack trace for devs (optional)"`
}

// TokenResponse represents a token exchange response
// @Description Returned on successful login or refresh
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOi..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOi..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"300"`
}

// EmailVerificationResponse contains the result of an email verification
// @Description Returned when email is verified and account activated
type EmailVerificationResponse struct {
	Verified bool   `json:"verified" example:"true"`
	Username string `json:"username" example:"johndoe"`
	Message  string `json:"message" example:"Email verified and Keycloak account created."`
}
