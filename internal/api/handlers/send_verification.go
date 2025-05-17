package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// SendVerificationLinkHandler godoc
// @Summary Resend verification email
// @Description Sends a new verification email to the user
// @Tags Email
// @Accept json
// @Produce plain
// @Param emailRequest body map[string]string true "Email payload"
// @Success 200 {string} string "Verification email sent"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "User not found"
// @Router /api/v1/auth/send-verification [post]
func SendVerificationLinkHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		fmt.Println("⚠️ Invalid request payload")
		return
	}

	if !isValidEmail(req.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		fmt.Printf("⚠️ Invalid email format: %s\n", req.Email)
		return
	}

	user, err := sqlite.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		fmt.Printf("❌ User not found: %s (err: %v)\n", req.Email, err)
		return
	}

	go func() {
		if err := SendVerificationEmail(user.Username, req.Email); err != nil {
			fmt.Printf("[!] Failed to resend email to %s: %v\n", user.Username, err)
		} else {
			fmt.Printf("📨 Verification email sent to: %s\n", req.Email)
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Verification email sent"))
}
