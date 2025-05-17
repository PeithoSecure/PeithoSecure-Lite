package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/peithosecure/peitho-backend/internal/auth/keycloak"
	corestub "github.com/peithosecure/peitho-backend/internal/corestub"
	"github.com/peithosecure/peitho-backend/internal/db/sqlite"
)

// DeleteAccountHandler godoc
// @Summary Delete a user account (admin only)
// @Description Deletes the specified user from Keycloak and logs the action in the audit trail.
// @Tags Admin
// @Accept json
// @Produce json
// @Param username query string true "Username to delete"
// @Success 200 {object} DeleteResponse "Account deleted successfully"
// @Failure 400 {object} map[string]string "Username missing or invalid"
// @Failure 500 {object} map[string]string "Failed to delete from Keycloak or log audit"
// @Router /api/v1/auth/delete [delete]
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		corestub.RespondWithTraceError(w, "username_required", http.StatusBadRequest)
		return
	}

	if err := keycloak.DeleteUser(GlobalConfig, username); err != nil {
		corestub.RespondWithTraceError(w, "account_delete_failed", http.StatusInternalServerError)
		return
	}

	_ = sqlite.LogAuditEvent(username, "account_deleted", r.RemoteAddr, r.UserAgent())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeleteResponse{
		Message: fmt.Sprintf("User '%s' deleted successfully", username),
	})
}

// DeleteResponse is returned when a user is successfully deleted
type DeleteResponse struct {
	Message string `json:"message" example:"User 'admin' deleted successfully"`
}
