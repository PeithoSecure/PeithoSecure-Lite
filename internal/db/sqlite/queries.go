package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/peithosecure/peitho-backend/internal/db/models"
)

// --- User queries ---

func GetUserByEmail(email string) (*models.User, error) {
	row := DB.QueryRow(`
		SELECT id, username, email, role, email_verified 
		FROM users 
		WHERE email = ?
	`, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.EmailVerified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

func GetUserByUsername(username string) (*models.User, error) {
	row := DB.QueryRow(`
		SELECT id, username, email, role, email_verified 
		FROM users 
		WHERE username = ?
	`, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.EmailVerified)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &user, err
}

// --- Email token logic ---

func CreateVerificationToken(username, token string, expiresAt time.Time) error {
	user, err := GetUserByUsername(username)
	if err != nil || user == nil {
		return errors.New("user not found")
	}
	return InsertEmailToken(user.Email, token, "verify")
}

func InsertEmailToken(email, token, tokenType string) error {
	_, err := GetDB().Exec(`
		INSERT INTO email_tokens (email, token, type, created_at)
		VALUES (?, ?, ?, datetime('now'))
	`, email, token, tokenType)
	return err
}

func GetUsernameByToken(token string) (string, error) {
	return GetUsernameByTokenAndType(token, "verify")
}

func GetUsernameByTokenAndType(token, tokenType string) (string, error) {
	var email string
	query := `
		SELECT email FROM email_tokens
		WHERE token = ? AND type = ?
		AND datetime(created_at, '+1 hour') > datetime('now')
	`

	err := GetDB().QueryRow(query, token, tokenType).Scan(&email)

	if err != nil {
		fmt.Printf("⚠️ Token lookup failed (%s): %v\n", tokenType, err)
	} else {
		fmt.Printf("✅ Token lookup success for: %s (type=%s)\n", email, tokenType)
	}
	fmt.Println("🕒 Backend UTC time:", time.Now().UTC().Format(time.RFC3339))

	return email, err
}

func DeleteVerificationToken(token string) error {
	_, err := GetDB().Exec(`DELETE FROM email_tokens WHERE token = ?`, token)
	return err
}

func MarkEmailVerified(username string) error {
	_, err := GetDB().Exec(`
		UPDATE users SET email_verified = 1 WHERE username = ?
	`, username)
	return err
}

// --- Audit Logging (Unified) ---

func LogAuditEvent(username, eventType, ip, userAgent string) error {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err := GetDB().Exec(`
		INSERT INTO audit_events (username, event_type, ip_address, user_agent, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, username, eventType, ip, userAgent, now)
	return err
}

func GetAuditEventsByUsername(username string, limit int) ([]models.AuditEvent, error) {
	rows, err := DB.Query(`
		SELECT id, username, event_type, ip_address, user_agent, created_at
		FROM audit_events
		WHERE username = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.AuditEvent
	for rows.Next() {
		var e models.AuditEvent
		var ts string
		if err := rows.Scan(&e.ID, &e.Username, &e.EventType, &e.IPAddress, &e.UserAgent, &ts); err != nil {
			return nil, err
		}
		e.CreatedAt, _ = time.ParseInLocation("2006-01-02 15:04:05", ts, time.UTC)
		events = append(events, e)
	}
	return events, nil
}
